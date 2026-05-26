package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func initSentry() error {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	sentryEnabled := os.Getenv("SENTRY_DSN") != ""
	if sentryEnabled {
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	if sentryEnabled {
		router.GET("/debug/sentry", func(c *gin.Context) {
			if hub := sentrygin.GetHubFromContext(c); hub != nil {
				hub.CaptureMessage("test sentry message")
			}
			panic("test sentry error")
		})
	}

	return router
}

func main() {
	if err := initSentry(); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := setupRouter().Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
