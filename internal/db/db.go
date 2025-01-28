package db

import (
	"context"
	"database/sql"
)

//go:generate mockgen -destination=mocks/row_mock.go -package=mocks . Row
type Row interface {
	Err() error
	Scan(dest ...any) error
}

//go:generate mockgen -destination=mocks/querier_mock.go -package=mocks . Querier
type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) Row
}

type repo struct {
	db *sql.DB
}

var _ Querier = (*repo)(nil)

func (r *repo) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return r.db.QueryRowContext(ctx, query, args...)
}
