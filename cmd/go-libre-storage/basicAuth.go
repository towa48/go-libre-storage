package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

const Realm = "Authorization Required"

func WebDavBasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		realmHeader := "Basic realm=" + strconv.Quote(Realm)
		header := ctx.Request.Header.Get("Authorization")
		user, found := searchCredential(header)
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			ctx.Header("WWW-Authenticate", realmHeader)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// set user for gin context
		ctx.Set(gin.AuthUserKey, user)
	}
}

func searchCredential(authHeader string) (string, bool) {
	if authHeader == "" {
		return "", false
	}
	headerParts := strings.Split(authHeader, " ")
	if headerParts[0] != "Basic" || len(headerParts) < 2 {
		return "", false
	}

	authValue, err := base64.StdEncoding.DecodeString(headerParts[1])
	if err != nil {
		return "", false
	}

	authParts := strings.Split(string(authValue), ":")
	if len(authParts) != 2 {
		return "", false
	}

	result := users.IsCredentialsValid(authParts[0], authParts[1])
	return authParts[0], result
}
