package users

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

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

	dk := encodePassword(password, salt)
	return bytes.Equal(dbPassword, dk)
}

func AddUser(login string, password string) int {
	_, found := GetUserByLogin(login)
	if found {
		fmt.Println("User already exists.")
		return 0
	}

	salt := generateSalt()
	dk := encodePassword(password, salt)
	dbSalt := base64.StdEncoding.EncodeToString(salt)
	dbPassword := base64.StdEncoding.EncodeToString(dk)

	user := User{
		Login:          login,
		PasswordHash:   dbPassword,
		Salt:           dbSalt,
		CreatedDateUtc: time.Now(),
	}

	return addUser(user)
}

func encodePassword(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keylen, sha1.New)
}

func generateSalt() []byte {
	token := make([]byte, keylen)
	rand.Read(token)
	return token
}
