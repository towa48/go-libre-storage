package users

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

const iterations int = 4096
const keylen int = 32

func IsCredentialsValid(login string, password string) bool {
	user, found := GetUserByLogin(login)
	if !found {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return false
	}

	dbPassword, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return false
	}

	dk := pbkdf2.Key([]byte(password), salt, iterations, keylen, sha1.New)
	return bytes.Equal(dbPassword, dk)
}

func generateSalt() []byte {
	token := make([]byte, keylen)
	rand.Read(token)
	return token
}
