package main

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

const templatesPath = "./web/templates/"

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", templatesPath + "_layout.html", templatesPath + "login.html")
	//r.AddFromFiles("index", templatesPath + "_layout.html", templatesPath + "index.html")

	return r
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.HTMLRender = createRender()

	Home(router)
	Accounts(router)
	WebDav(router)

	return router
}

func main() {
	router := setupRouter()
	router.Run(":3000")
}
