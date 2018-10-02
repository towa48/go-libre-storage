package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {

		viewModel := gin.H{
			"title": "Login",
		}

		c.HTML(http.StatusOK, "login", viewModel)
	})
}