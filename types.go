package goopy

import (
	"context"
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
