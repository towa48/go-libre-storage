package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/towa48/go-libre-storage/internal/pkg/assets"
)

func Home(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {

		manifest, _ := assets.GetAssetsManifest()
		viewModel := gin.H{
			"title":    "Go Libre Storage",
			"chunkUrl": manifest.ChunkScriptUrl,
		}

		if isAuthenticated(c) {
			runtimeContent, _ := assets.GetAssetContent(manifest.MainRuntimeUrl)
			viewModel["runtimeContent"] = runtimeContent
			viewModel["styleUrl"] = manifest.MainStyleUrl
			viewModel["scriptUrl"] = manifest.MainScriptUrl

			c.HTML(http.StatusOK, "index", viewModel)
			return
		}

		runtimeContent, _ := assets.GetAssetContent(manifest.WelcomeRuntimeUrl)
		viewModel["runtimeContent"] = runtimeContent
		viewModel["styleUrl"] = manifest.WelcomeStyleUrl
		viewModel["scriptUrl"] = manifest.WelcomeScriptUrl

		c.HTML(http.StatusOK, "welcome", viewModel)
	})
}
