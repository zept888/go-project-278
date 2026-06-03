package store

import (
	"context"

	"github.com/zept888/go-project-278/internal/db"
)

func (p *Postgres) CountVisits(ctx context.Context) (int64, error) {
	return p.q.CountLinkVisits(ctx)
}

func (p *Postgres) ListVisits(ctx context.Context, offset, limit int) ([]LinkVisit, error) {
	if limit <= 0 {
		return []LinkVisit{}, nil
	}
	rows, err := p.q.ListLinkVisitsPage(ctx, db.ListLinkVisitsPageParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}
	out := make([]LinkVisit, len(rows))
	for i, row := range rows {
		out[i] = toVisit(row)
	}
	return out, nil
}

func (p *Postgres) CreateVisit(ctx context.Context, params CreateVisitParams) (LinkVisit, error) {
	row, err := p.q.CreateLinkVisit(ctx, db.CreateLinkVisitParams{
		LinkID:    params.LinkID,
		Ip:        params.IP,
		UserAgent: params.UserAgent,
		Referer:   params.Referer,
		Status:    int32(params.Status),
	})
	if err != nil {
		return LinkVisit{}, err
	}
	return toVisit(row), nil
}

func toVisit(row db.LinkVisit) LinkVisit {
	return LinkVisit{
		ID:        row.ID,
		LinkID:    row.LinkID,
		IP:        row.Ip,
		UserAgent: row.UserAgent,
		Referer:   row.Referer,
		Status:    int(row.Status),
		CreatedAt: row.CreatedAt.Time,
	}
}
