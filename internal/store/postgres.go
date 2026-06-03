package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zept888/go-project-278/internal/db"
)

type Postgres struct {
	q *db.Queries
}

func NewPostgres(q *db.Queries) *Postgres {
	return &Postgres{q: q}
}

func (p *Postgres) List(ctx context.Context) ([]Link, error) {
	rows, err := p.q.ListLinks(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Link, len(rows))
	for i, row := range rows {
		out[i] = toLink(row)
	}
	return out, nil
}

func (p *Postgres) Get(ctx context.Context, id int64) (Link, error) {
	row, err := p.q.GetLink(ctx, id)
	if err != nil {
		return Link{}, mapErr(err)
	}
	return toLink(row), nil
}

func (p *Postgres) GetByShortName(ctx context.Context, shortName string) (Link, error) {
	row, err := p.q.GetLinkByShortName(ctx, shortName)
	if err != nil {
		return Link{}, mapErr(err)
	}
	return toLink(row), nil
}

func (p *Postgres) Create(ctx context.Context, originalURL, shortName string) (Link, error) {
	row, err := p.q.CreateLink(ctx, db.CreateLinkParams{
		OriginalUrl: originalURL,
		ShortName:   shortName,
	})
	if err != nil {
		return Link{}, mapErr(err)
	}
	return toLink(row), nil
}

func (p *Postgres) Update(ctx context.Context, id int64, originalURL, shortName string) (Link, error) {
	row, err := p.q.UpdateLink(ctx, db.UpdateLinkParams{
		ID:          id,
		OriginalUrl: originalURL,
		ShortName:   shortName,
	})
	if err != nil {
		return Link{}, mapErr(err)
	}
	return toLink(row), nil
}

func (p *Postgres) Delete(ctx context.Context, id int64) error {
	if _, err := p.q.GetLink(ctx, id); err != nil {
		return mapErr(err)
	}
	return p.q.DeleteLink(ctx, id)
}

func toLink(row db.Link) Link {
	return Link{
		ID:          row.ID,
		OriginalURL: row.OriginalUrl,
		ShortName:   row.ShortName,
	}
}

func mapErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return err
}
