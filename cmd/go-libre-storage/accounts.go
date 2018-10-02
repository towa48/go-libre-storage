package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Accounts(route *gin.Engine) {
	accounts := route.Group("/accounts")

	accounts.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "Login")
	})
}