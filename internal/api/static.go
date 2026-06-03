package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterStatic(router *gin.Engine) {
	dir := os.Getenv("STATIC_DIR")
	if dir == "" {
		dir = "/app/public"
	}
	index := filepath.Join(dir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return
	}

	log.Printf("serving UI static files from %s", dir)

	assets := filepath.Join(dir, "assets")
	if _, err := os.Stat(assets); err == nil {
		router.Static("/assets", assets)
	}

	serveSPA := func(c *gin.Context) {
		c.File(index)
	}
	router.GET("/", serveSPA)
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/r/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Status(http.StatusNotFound)
			return
		}
		serveSPA(c)
	})
}
