package main

import (
	"github.com/towa48/go-libre-storage/internal/pkg/config"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

const templatesPath = "./web/templates/"

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("welcome", templatesPath+"_layout.html", templatesPath+"welcome.html")
	r.AddFromFiles("index", templatesPath+"_layout.html", templatesPath+"index.html")

	return r
}

func createSessionStore() gin.HandlerFunc {
	store := cookie.NewStore([]byte(config.Get().CookieSecret))
	return sessions.Sessions("libre_session", store)
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(createSessionStore())
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.HTMLRender = createRender()

	// static
	router.Static("/js", "./web/wwwroot/js")
	router.Static("/styles", "./web/wwwroot/styles")
	router.Static("/images", "./web/wwwroot/images")
	router.Static("/fonts", "./web/wwwroot/fonts")
	router.StaticFile("/robots.txt", "./web/wwwroot/robots.txt")
	router.StaticFile("/site.webmanifest", "./web/wwwroot/site.webmanifest")
	router.StaticFile("/favicon.ico", "./web/wwwroot/favicon.ico")
	router.StaticFile("/icon.png", "./web/wwwroot/icon.png")

	Home(router)
	Accounts(router)
	WebDav(router)

	return router
}

func main() {
	users.CheckDatabase()
	files.CheckDatabase()

	if checkCliCommands() {
		return
	}

	router := setupRouter()
	router.Run(":3000")
}
