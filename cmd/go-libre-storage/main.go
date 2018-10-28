package main

import (
	"os"

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

	args := os.Args[1:]
	crawlArg := contains(args, "--crawl")

	if crawlArg {
		crawl()
		return
	}

	router := setupRouter()
	router.Run(":3000")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
