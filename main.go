package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/zept888/go-project-278/internal/api"
	"github.com/zept888/go-project-278/internal/db"
	"github.com/zept888/go-project-278/internal/linkutil"
	"github.com/zept888/go-project-278/internal/store"
)

func initSentry() error {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return nil
	}
	return sentry.Init(sentry.ClientOptions{Dsn: dsn})
}

func setupRouter(st store.Store, baseURL string) *gin.Engine {
	router := gin.Default()
	router.Use(api.CORSMiddleware())

	sentryEnabled := os.Getenv("SENTRY_DSN") != ""
	if sentryEnabled {
		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
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

	if st != nil {
		h := &api.Links{Store: st, BaseURL: baseURL}
		router.GET("/r/:shortName", h.Redirect)
		links := router.Group("/api/links")
		links.GET("", h.List)
		links.POST("", h.Create)
		links.GET("/:id", h.Get)
		links.PUT("/:id", h.Update)
		links.DELETE("/:id", h.Delete)
	}

	api.RegisterStatic(router)

	return router
}

func openStore(ctx context.Context) (store.Store, func(), error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL is not set, using in-memory store (local dev only)")
		return store.NewMemory(), func() {}, nil
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, func() {}, err
	}
	return store.NewPostgres(db.New(pool)), pool.Close, nil
}

func baseURL(port string) string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return linkutil.NormalizeBaseURL(v)
	}
	return "http://localhost:" + port
}

func main() {
	_ = godotenv.Load()

	if err := initSentry(); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	st, closeDB, err := openStore(ctx)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer closeDB()

	if err := setupRouter(st, baseURL(port)).Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
