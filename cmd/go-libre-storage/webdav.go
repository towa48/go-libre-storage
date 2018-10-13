package main

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
)

const Prefix string = "/webdav"
const XmlDocumentType string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
const WebDavStatusOk string = "HTTP/1.1 200 OK"

func WebDav(r *gin.Engine) {
	authorized := r.Group(Prefix, WebDavBasicAuth())

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
		c.Status(403)
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
		payload := files.GetPathInfo(path, 0, includeContent)

		resp := getMultistatusResponse(payload)

		httpStatus := http.StatusMultiStatus
		c.Status(httpStatus)
		c.Writer.Write([]byte(XmlDocumentType))
		c.XML(httpStatus, resp)
	})
}

func stripPrefix(path string) string {
	if result := strings.TrimPrefix(path, Prefix); len(result) < len(path) {
		return result
	}

	return path
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

func getMultistatusResponse(payload []files.FileInfo) Multistatus {
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
						DateTag:          fi.ETag,
						ContentType:      fi.ContentType,
						ContentLength:    fi.ContentLength,
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
	DateTag          string      `xml:"d:prop>d:getetag"`
	CreationDate     string      `xml:"d:prop>d:creationdate"`
	DisplayName      string      `xml:"d:prop>d:displayname"`
	LastModifiedDate string      `xml:"d:prop>d:getlastmodified"`
	ContentType      string      `xml:"d:prop>d:getcontenttype"`
	ContentLength    string      `xml:"d:prop>d:getcontentlength"`
	ResourceType     interface{} `xml:"d:prop>d:resourcetype"`
}

type Property struct {
	XMLName xml.Name
}

type CollectionResourceType struct {
	XMLName   xml.Name `xml:"d:resourcetype"`
	FakeValue string   `xml:"d:collection"`
}
