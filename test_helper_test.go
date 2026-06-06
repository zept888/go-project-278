package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/zept888/go-project-278/internal/db"
	"github.com/zept888/go-project-278/internal/store"
)

func testDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"
}

func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Setenv("SENTRY_DSN", "")
	return setupRouter(testStore(t), testBaseURL)
}

func testStore(t *testing.T) store.Store {
	t.Helper()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, testDSN())
	if err != nil {
		t.Fatalf("connect postgres: %v (run: docker compose up -d)", err)
	}
	t.Cleanup(pool.Close)

	if err := migrateTestDB(pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := resetTestDB(ctx, pool); err != nil {
		t.Fatalf("reset db: %v", err)
	}

	return store.NewPostgres(db.New(pool))
}

func migrateTestDB(pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Join(err, sqlDB.Close())
	}
	if err := goose.Up(sqlDB, "db/migrations"); err != nil {
		return errors.Join(err, sqlDB.Close())
	}
	return sqlDB.Close()
}

func resetTestDB(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "TRUNCATE link_visits, links RESTART IDENTITY CASCADE")
	return err
}
