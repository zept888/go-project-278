package store

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("link not found")
	ErrConflict = errors.New("short_name already exists")
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortName   string
}

type LinkVisit struct {
	ID        int64
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
	Status    int
	CreatedAt time.Time
}

type CreateVisitParams struct {
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
	Status    int
}

type Store interface {
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context, offset, limit int) ([]Link, error)
	Get(ctx context.Context, id int64) (Link, error)
	GetByShortName(ctx context.Context, shortName string) (Link, error)
	Create(ctx context.Context, originalURL, shortName string) (Link, error)
	Update(ctx context.Context, id int64, originalURL, shortName string) (Link, error)
	Delete(ctx context.Context, id int64) error

	CountVisits(ctx context.Context) (int64, error)
	ListVisits(ctx context.Context, offset, limit int) ([]LinkVisit, error)
	CreateVisit(ctx context.Context, params CreateVisitParams) (LinkVisit, error)
}
