package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const Prefix string = "/webdav"

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
		fmt.Println("Request: " + path)

		c.Status(403)
	})
}

func stripPrefix(path string) string {
	if result := strings.TrimPrefix(path, Prefix); len(result) < len(path) {
		return result
	}

	return path
}
