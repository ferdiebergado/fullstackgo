package db

import (
	"context"
)

type MockRow struct {
	ScanFn func(dest ...any) error
	ErrFn  func() error
}

func (m *MockRow) Err() error {
	if m.ErrFn != nil {
		return m.ErrFn()
	}

	return nil
}

func (m *MockRow) Scan(dest ...any) error {
	if m.ScanFn != nil {
		return m.ScanFn(dest...)
	}
	return nil
}

var _ Row = (*MockRow)(nil)

type MockDB struct {
	QueryRowContextFn func(tx context.Context, query string, args ...any) Row
}

var _ Querier = (*MockDB)(nil)

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	if m.QueryRowContextFn != nil {
		return m.QueryRowContextFn(ctx, query, args...)
	}

	return nil
}
