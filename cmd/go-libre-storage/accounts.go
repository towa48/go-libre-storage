package main

import (
	"net/http"

	"github.com/towa48/go-libre-storage/internal/pkg/users"

	"github.com/gin-gonic/gin"
)

func Accounts(route *gin.Engine) {
	accounts := route.Group("/accounts")

	accounts.POST("/login", func(c *gin.Context) {
		var loginParams LoginParams
		c.BindJSON(&loginParams)

		if loginParams.Username == EmptyString || loginParams.Password == EmptyString {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Params"})
			return
		}

		isValid := users.IsCredentialsValid(loginParams.Username, loginParams.Password)
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "InvalidCredentials"})
			return
		}

		authenticate(c, loginParams.Username, loginParams.RememberMe)
		c.JSON(http.StatusOK, gin.H{"Url": "/"})
	})
}

type LoginParams struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}
