package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebDav(r *gin.Engine) {
	authorized := r.Group("/webdav", WebDavBasicAuth())

	authorized.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WebDav")
	})
}
