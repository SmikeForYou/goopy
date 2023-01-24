package goopy

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
)

type ViewModel interface {
	Table() string
}

type Model interface {
	ViewModel
	Pk() (string, any)
}

type rows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type exer interface {
	pgxscan.Querier
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type querier[T Model] interface {
	Select(ctx context.Context, expressions ...goqu.Expression) ([]T, error)
	SelectOne(ctx context.Context, expressions ...goqu.Expression) (T, error)
}

type manager[T Model] interface {
	Insert(ctx context.Context, in T) (T, error)
	BulkInsert(ctx context.Context, in ...T) ([]T, error)
	Update(ctx context.Context, in T) (T, error)
	Delete(ctx context.Context, in T) error
}
