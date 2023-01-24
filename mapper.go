package goopy

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var PostgresDialect = goqu.Dialect("postgres")

type Querier[T ViewModel] struct {
	conn      pgxscan.Querier
	tableName string
	dialect   string
}

func NewPGQuerier[T ViewModel](conn pgxscan.Querier) *Querier[T] {
	var t T
	return &Querier[T]{
		conn:      conn,
		tableName: t.Table(),
		dialect:   "postgres",
	}
}

func (q Querier[T]) selectMany(ctx context.Context, sql string, args []any) ([]T, error) {
	target := make([]T, 0)
	err := pgxscan.Select(ctx, q.conn, &target, sql, args...)
	return target, err
}

func (q Querier[T]) selectOne(ctx context.Context, sql string, args []any) (T, error) {
	var target T
	err := pgxscan.Get(ctx, q.conn, &target, sql, args...)
	return target, err
}

func (q Querier[T]) Select(ctx context.Context, expressions ...goqu.Expression) ([]T, error) {
	builder := toSelectQuery(q.tableName, expressions...).WithDialect(q.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		return nil, NewError(builder, err)
	}
	res, err := q.selectMany(ctx, query, args)
	if err != nil {
		return nil, NewError(builder, err)
	}
	return res, nil
}

func (q Querier[T]) SelectOne(ctx context.Context, expressions ...goqu.Expression) (T, error) {
	builder := toSelectQuery(q.tableName, expressions...).WithDialect(q.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		var t T
		return t, NewError(builder, err)
	}
	res, err := q.selectOne(ctx, query, args)
	if err != nil {
		return res, NewError(builder, err)
	}
	return res, nil
}

type Manager[T Model] struct {
	*Querier[T]
	conn exer
}

func NewManager[T Model](conn exer) *Manager[T] {
	return &Manager[T]{
		Querier: NewPGQuerier[T](conn),
		conn:    conn,
	}
}

func (m Manager[T]) exec(ctx context.Context, query string, args []any) (int64, error) {
	tag, err := m.conn.Exec(ctx, query, args...)
	return tag.RowsAffected(), err
}

func (m Manager[T]) Insert(ctx context.Context, in T) (T, error) {
	builder := toInsertQuery(m.tableName, in).WithDialect(m.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		var t T
		return t, NewError(builder, err)
	}
	res, err := m.selectOne(ctx, query, args)
	if err != nil {
		return res, NewError(builder, err)
	}
	return res, nil
}

func (m Manager[T]) BulkInsert(ctx context.Context, in ...T) ([]T, error) {
	builder := toInsertQuery(m.tableName, in).WithDialect(m.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		return nil, NewError(builder, err)
	}
	res, err := m.selectMany(ctx, query, args)
	if err != nil {
		return res, NewError(builder, err)
	}
	return res, nil
}

func (m Manager[T]) Update(ctx context.Context, in T) (T, error) {
	builder := toUpdateQuery(m.tableName, in).WithDialect(m.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		var t T
		return t, NewError(builder, err)
	}
	res, err := m.selectOne(ctx, query, args)
	if err != nil {
		return res, NewError(builder, err)
	}
	return res, nil
}

func (m Manager[T]) Delete(ctx context.Context, in T) error {
	builder := toDeleteQuery(m.tableName, in).WithDialect(m.dialect)
	query, args, err := builder.Prepared(true).ToSQL()
	if err != nil {
		return NewError(builder, err)
	}
	_, err = m.exec(ctx, query, args)
	if err != nil {
		return NewError(builder, err)
	}
	return err
}
