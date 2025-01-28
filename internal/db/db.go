//go:generate mockgen -source=db.go -destination=mocks/db_mock.go -package=db_mock
package db

import (
	"context"
	"database/sql"
)

type Row interface {
	Err() error
	Scan(dest ...any) error
}

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
