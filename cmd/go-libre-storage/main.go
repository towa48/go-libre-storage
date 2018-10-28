package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

const templatesPath = "./web/templates/"

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", templatesPath+"_layout.html", templatesPath+"login.html")
	//r.AddFromFiles("index", templatesPath + "_layout.html", templatesPath + "index.html")

	return r
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.HTMLRender = createRender()

	Home(router)
	Accounts(router)
	WebDav(router)

	return router
}

func main() {
	users.CheckDatabase()
	files.CheckDatabase()

	doCrawl := flag.Bool("crawl", false, "Clear DB file metadata and restore it from filesystem. All shared items will be dropped.")
	newUser := flag.String("adduser", "", "Add new user. Password will be promted.")
	flag.Parse()

	if doCrawl != nil && *doCrawl {
		crawl()
		return
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
		return
	}

	router := setupRouter()
	router.Run(":3000")
}

func trim(s string) string {
	return strings.Trim(s, " \n")
}

func isEmptyString(s string) bool {
	return trim(s) == ""
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
