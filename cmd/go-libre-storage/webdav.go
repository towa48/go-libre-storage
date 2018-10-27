package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/towa48/go-libre-storage/internal/pkg/config"

	"github.com/towa48/go-libre-storage/internal/pkg/users"

	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
)

const WebDavPrefix string = "/webdav"
const XmlDocumentType string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
const WebDavStatusOk string = "HTTP/1.1 200 OK"

func WebDav(r *gin.Engine) {
	authorized := r.Group(WebDavPrefix, WebDavBasicAuth())

	authorized.OPTIONS("/*path", func(c *gin.Context) {
		path := stripPrefix(c.Request.URL.Path)
		if path == "/" {
			c.Header("Allow", "OPTIONS, PROPFIND")
		} else {
			// TBD
		}

		c.Header("DAV", "1, 2")
		c.Header("MS-Author-Via", "DAV")
	})

	authorized.Handle("GET", "/*path", func(c *gin.Context) {
		path := stripPrefix(c.Request.URL.Path)

		login := c.MustGet(gin.AuthUserKey).(string)
		user, found := users.GetUserByLogin(login)
		if !found {
			c.Status(http.StatusForbidden)
			return
		}

		fi, hasAccess := files.GetFileInfo(path, user.Id, WebDavPrefix)
		if !hasAccess {
			c.Status(http.StatusForbidden)
			return
		}
		if fi.IsDir {
			c.Status(http.StatusBadRequest)
			return
		}

		filePathRoot, found := files.GetFileHierarchy(fi.Id)
		if !found {
			fmt.Printf("File %d path not found\n", fi.Id)
			c.String(http.StatusInternalServerError, "Internal server error")
			return
		}

		filePath := buildFilePath(filePathRoot)
		fmt.Println(filePath)

		file, err := os.Open(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "%v", err)
			return
		}
		defer file.Close()

		c.Header("Content-Type", fi.Mime)
		c.Header("ETag", fi.ETag)
		c.Header("Last-Modified", fi.ModifiedDateUtc.Format(time.RFC1123))

		io.Copy(c.Writer, file)
	})

	authorized.Handle("PROPFIND", "/*path", func(c *gin.Context) {
		path := stripPrefix(c.Request.URL.Path)
		depth := parseDepth(c.Request.Header.Get("Depth"))

		if depth == invalidDepth || depth == infiniteDepth {
			c.Status(http.StatusBadRequest)
			return
		}

		data, err := c.GetRawData()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		var req Propfind
		err = xml.Unmarshal(data, &req)
		// TODO: analyze request payload

		includeContent := depth == 1
		login := c.MustGet(gin.AuthUserKey).(string)

		user, found := users.GetUserByLogin(login)
		if !found {
			c.Status(http.StatusForbidden)
			return
		}

		payload, hasAccess := files.GetFolderContent(path, user.Id, WebDavPrefix, includeContent)
		if !hasAccess {
			c.Status(http.StatusForbidden)
			return
		}
		if hasAccess && payload == nil {
			c.String(http.StatusNotFound, "text/plain", "Folder was deleted.")
			return
		}

		resp := getMultistatusResponse(payload)

		httpStatus := http.StatusMultiStatus
		c.Status(httpStatus)
		c.Writer.Write([]byte(XmlDocumentType))
		c.XML(httpStatus, resp)
	})

	authorized.Handle("MKCOL", "/*path", func(c *gin.Context) {
		path := stripPrefix(c.Request.URL.Path)

		login := c.MustGet(gin.AuthUserKey).(string)
		user, found := users.GetUserByLogin(login)
		if !found {
			c.Status(http.StatusForbidden)
			return
		}

		dir := getFileSystemPath(path, user)

		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		root := parseUrl(path)
		eRoot, parentId, found := getFirstUnknownFolder(root, user.Id)

		if !found {
			c.Status(http.StatusBadRequest)
			return
		}

		db := files.GetDbConnection()
		defer db.Close()

		t := time.Now()
		for eRoot != nil {
			fi := files.DbFileInfo{
				Name:            eRoot.Name,
				Path:            eRoot.Url,
				CreatedDateUtc:  t,
				ModifiedDateUtc: t,
				OwnerId:         user.Id,
			}
			parentId = files.AppendFolder(db, fi, parentId)
			eRoot = eRoot.Child
		}

		c.String(http.StatusCreated, "")
	})

	authorized.DELETE("/*path", func(c *gin.Context) {
		path := stripPrefix(c.Request.URL.Path)

		login := c.MustGet(gin.AuthUserKey).(string)
		user, found := users.GetUserByLogin(login)
		if !found {
			c.Status(http.StatusForbidden)
			return
		}

		isFolder := strings.HasSuffix(path, UrlSeparator)

		if isFolder {
			fi, found := files.GetFolderInfo(path, user.Id)
			if !found {
				c.String(http.StatusBadRequest, "")
				return
			}

			fsp := getFileSystemPath(path, user)
			err := os.RemoveAll(fsp)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			files.RemoveFolder(fi.Id)
			c.Status(http.StatusNoContent)
			return
		}

		fi, found := files.GetFileInfo(path, user.Id, WebDavPrefix)
		if !found {
			c.String(http.StatusBadRequest, "")
			return
		}

		fsp := getFileSystemPath(path, user)
		err := os.RemoveAll(fsp)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		files.RemoveFile(fi.Id)
		c.Status(http.StatusNoContent)
	})
}

func stripPrefix(path string) string {
	if result := strings.TrimPrefix(path, WebDavPrefix); len(result) < len(path) {
		if !strings.HasPrefix(result, "/") {
			return "/" + result
		}

		return result
	}

	return path
}

func buildFilePath(root files.DbHierarchyItem) string {
	separator := string(os.PathSeparator)

	result := config.Get().Storage
	if !strings.HasSuffix(result, separator) {
		result = result + separator
	}

	item := root
	for item.Child != nil {
		result = result + item.Name + separator
		item = *item.Child
	}

	return result + item.Name
}

func getFileSystemPath(url string, user users.User) string {
	separator := string(os.PathSeparator)
	r := config.Get().Storage

	if !strings.HasSuffix(r, separator) {
		r = r + separator
	}

	if !strings.HasPrefix(url, UrlSeparator) {
		url = UrlSeparator + url
	}

	return r + user.Login + url
}

func parseUrl(url string) *UrlHierarchyItem {
	p := url
	i := strings.Index(p, UrlSeparator)

	root := &UrlHierarchyItem{
		Url:   "/",
		IsDir: true,
	}
	prev := root

	for i != -1 {
		if i == 0 {
			p = strings.TrimPrefix(p, UrlSeparator)
		}

		i = strings.Index(p, UrlSeparator)
		var next *UrlHierarchyItem
		if i > 0 {
			name := p[0:i]
			next = &UrlHierarchyItem{
				Name:  name,
				Url:   prev.Url + name + UrlSeparator,
				IsDir: true,
			}
			p = strings.TrimPrefix(p, name)
		} else if i == -1 && len(p) > 0 {
			next = &UrlHierarchyItem{
				Name: p,
				Url:  prev.Url + p,
			}
		}

		if next != nil {
			prev.Child = next
			prev = next
		}
	}

	return root
}

func getFirstUnknownFolder(root *UrlHierarchyItem, userId int) (folder *UrlHierarchyItem, parentId int64, found bool) {
	var pid int64
	var result bool
	for root != nil {
		fi, fo := files.GetFolderInfo(root.Url, userId)
		if !fo {
			break
		}

		result = true
		pid = fi.Id
		root = root.Child
	}

	if !result {
		return nil, pid, result
	}

	return root, pid, result
}

const (
	infiniteDepth = -1
	invalidDepth  = -2
)

func parseDepth(s string) int {
	switch s {
	case "0":
		return 0
	case "1":
		return 1
	case "infinity":
		return infiniteDepth
	}
	return invalidDepth
}

func getMultistatusResponse(payload []files.DbFileInfo) Multistatus {
	var responses []MultistatusResponse

	for _, fi := range payload {
		if fi.IsDir {
			responses = append(responses, MultistatusResponse{
				Href: fi.Path,
				Props: []interface{}{
					DirPropStat{
						Status:           WebDavStatusOk,
						DisplayName:      fi.Name,
						CreationDate:     fi.CreatedDateUtc.Format(time.RFC3339),
						LastModifiedDate: fi.ModifiedDateUtc.Format(time.RFC1123),
						ResourceType:     &CollectionResourceType{},
					},
				},
			})
		} else {
			responses = append(responses, MultistatusResponse{
				Href: fi.Path,
				Props: []interface{}{
					FilePropStat{
						Status:           WebDavStatusOk,
						DisplayName:      fi.Name,
						CreationDate:     fi.CreatedDateUtc.Format(time.RFC3339),
						LastModifiedDate: fi.ModifiedDateUtc.Format(time.RFC1123),
						ETag:             fi.ETag,
						ContentType:      fi.Mime,
						ContentLength:    strconv.FormatInt(fi.Size, 10),
					},
				},
			})
		}
	}

	return Multistatus{
		XmlNs:   "DAV:",
		Reponse: responses,
	}
}

type Propfind struct {
	XmlName xml.Name `xml:"propfind"`
	XmlNs   string   `xml:"xmlns,attr"`
	Allprop string   `xml:"allprop"`
}

type Multistatus struct {
	XMLName xml.Name              `xml:"d:multistatus"`
	XmlNs   string                `xml:"xmlns:d,attr"`
	Reponse []MultistatusResponse `xml:"d:response"`
}

type MultistatusResponse struct {
	Href  string        `xml:"d:href"`
	Props []interface{} `xml:"d:propstat"`
}

type DirPropStat struct {
	Status           string      `xml:"d:status"`
	CreationDate     string      `xml:"d:prop>d:creationdate"`
	DisplayName      string      `xml:"d:prop>d:displayname"`
	LastModifiedDate string      `xml:"d:prop>d:getlastmodified"`
	ResourceType     interface{} `xml:"d:prop>d:resourcetype"`
}

type FilePropStat struct {
	Status           string      `xml:"d:status"`
	ETag             string      `xml:"d:prop>d:getetag"`
	CreationDate     string      `xml:"d:prop>d:creationdate"`
	DisplayName      string      `xml:"d:prop>d:displayname"`
	LastModifiedDate string      `xml:"d:prop>d:getlastmodified"`
	ContentType      string      `xml:"d:prop>d:getcontenttype"`
	ContentLength    string      `xml:"d:prop>d:getcontentlength"`
	ResourceType     interface{} `xml:"d:prop>d:resourcetype"`
}

type CollectionResourceType struct {
	XMLName   xml.Name `xml:"d:resourcetype"`
	FakeValue string   `xml:"d:collection"`
}

type UrlHierarchyItem struct {
	Name  string
	IsDir bool
	Url   string
	Child *UrlHierarchyItem
}
