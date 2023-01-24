package goopy

import (
	"errors"
	"github.com/doug-martin/goqu/v9/exp"
)

var (
	ErrNoRows = errors.New("scanning one: no rows in result set")
)

func NewError(query exp.SQLExpression, err error) error {
	if err == nil {
		return nil
	}
	return Err{
		Query: query,
		Err:   err,
	}
}

type Err struct {
	Query exp.SQLExpression
	Err   error
}

func (e Err) Error() string {
	return e.Err.Error()
}

func (e Err) IsNotFound() bool {
	return e.Error() == ErrNoRows.Error()
}
