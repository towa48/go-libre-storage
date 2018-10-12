package main

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const Prefix string = "/webdav"
const XmlDocumentType string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"

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

		if path == "/" && depth == 0 {
			resp := getMultistatusResponse("/") // TODO: get root folder properties

			c.Writer.Write([]byte(XmlDocumentType))
			c.XML(http.StatusMultiStatus, resp)
			return
		}

		c.Status(403)
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

func getMultistatusResponse(url string) Multistatus {
	timeNow := time.Now()
	createdTimeVal := timeNow.Format(time.RFC3339)
	modifiedTimeVal := timeNow.Format(time.RFC1123)

	props := []PropStat{
		{
			Status:           "HTTP/1.1 200 OK",
			CreationDate:     createdTimeVal,
			DisplayName:      "box",
			LastModifiedDate: modifiedTimeVal,
			ResourceType:     &CollectionResourceType{},
		},
	}

	return Multistatus{
		XmlNs: "DAV:",
		Reponse: []MultistatusResponse{
			{
				Href:  "/",
				Props: props,
			},
		},
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
	Href  string     `xml:"d:href"`
	Props []PropStat `xml:"d:propstat"`
}

type PropStat struct {
	Status           string      `xml:"d:status"`
	DateTag          string      `xml:"d:prop>d:getetag"` // only files, etag?
	CreationDate     string      `xml:"d:prop>d:creationdate"`
	DisplayName      string      `xml:"d:prop>d:displayname"`
	LastModifiedDate string      `xml:"d:prop>d:getlastmodified"`
	ContentType      string      `xml:"d:prop>d:getcontenttype"`   // only files, mime
	ContentLength    string      `xml:"d:prop>d:getcontentlength"` // only files, bytes
	ResourceType     interface{} `xml:"d:prop>d:resourcetype"`
}

type CollectionResourceType struct {
	XMLName   xml.Name `xml:"d:resourcetype"`
	FakeValue string   `xml:"d:collection"`
}
