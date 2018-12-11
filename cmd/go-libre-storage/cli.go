package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

func checkCliCommands() bool {
	doCrawl := flag.Bool("crawl", false, "Clear DB file metadata and restore it from filesystem. All shared items will be dropped.")
	newUser := flag.String("adduser", "", "Add new user. Password will be promted.")
	flag.Parse()

	if doCrawl != nil && *doCrawl {
		crawl()
		return true
	}

	if newUser != nil && *newUser != "" {
		login := *newUser
		reader := bufio.NewReader(os.Stdin)
		pass := ""
		for isEmptyString(pass) {
			fmt.Print("Enter password: ")
			pass, _ = reader.ReadString('\n')
			if isEmptyString(pass) {
				fmt.Println("Password is invalid.")
			}
		}

		pass = trim(pass)
		userId := users.AddUser(login, pass)

		if userId != 0 {
			createUserRoot(userId, login)
		}
		return true
	}

	return false
}

func trim(s string) string {
	return strings.Trim(s, " \n")
}

func isEmptyString(s string) bool {
	return trim(s) == ""
}

/*func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}*/
