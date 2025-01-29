package db

import (
	"context"
	"database/sql"
	"errors"
)

var ErrDuplicateUser = errors.New("user already exists")
var ErrNullValue = errors.New("not null constraint violation")

type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type repo struct {
	db *sql.DB
}

var _ Querier = (*repo)(nil)

func (r *repo) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}
