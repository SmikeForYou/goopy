package goopy

import (
	"github.com/doug-martin/goqu/v9"
)

func toSelectQuery(tableName string, expressions ...goqu.Expression) *goqu.SelectDataset {
	return goqu.Select("*").From(tableName).Where(expressions...)
}

func toInsertQuery(tableName string, in ...any) *goqu.InsertDataset {
	return goqu.Insert(tableName).Rows(in...).Returning("*")
}

func toUpdateQuery[T Model](tableName string, in T) *goqu.UpdateDataset {
	pkcol, pkval := in.Pk()
	return goqu.Update(tableName).Set(in).Where(goqu.C(pkcol).Eq(pkval)).Returning("*")
}

func toDeleteQuery[T Model](tableName string, in T) *goqu.DeleteDataset {
	pkcol, pkval := in.Pk()
	return goqu.Delete(tableName).Where(goqu.C(pkcol).Eq(pkval))
}
