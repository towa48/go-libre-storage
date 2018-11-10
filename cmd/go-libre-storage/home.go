package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {

		viewModel := gin.H{
			"title": "Go Libre Storage",
		}

		c.HTML(http.StatusOK, "welcome", viewModel)
	})
}
