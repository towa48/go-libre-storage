package main

import (
	"fmt"
	"net/http"

	"github.com/towa48/go-libre-storage/internal/pkg/users"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const AuthCookie = "_gls_auth"

func cookieAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		val := session.Get(AuthCookie)

		if val == nil || val == EmptyString {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		login := val.(string)

		found := userExists(login)
		if !found {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// set user for gin context
		ctx.Set(gin.AuthUserKey, login)
	}
}

func authenticate(ctx *gin.Context, login string, rememberMe bool) {
	session := sessions.Default(ctx)
	session.Set(AuthCookie, login)
	err := session.Save()

	fmt.Println("Session save error.", err)
	// TODO: remember me
}

func slideExpiration(ctx *gin.Context) {
	// TODO
}

func isAuthenticated(ctx *gin.Context) bool {
	session := sessions.Default(ctx)
	val := session.Get(AuthCookie)

	if val == nil || val == EmptyString {
		return false
	}

	login := val.(string)
	return userExists(login)
}

func userExists(login string) bool {
	_, f := users.GetUserByLogin(login)
	return f
}
