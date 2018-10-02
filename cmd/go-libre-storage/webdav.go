package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WebDav(r *gin.Engine) {
	authorized := r.Group("/webdav", gin.BasicAuth(gin.Accounts{
        "user": "password",
    }))

	authorized.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "WebDav")
	})
}