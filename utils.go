package goopy

import (
	"github.com/doug-martin/goqu/v9"
)

func toSelectQuery(tableName string, expressions ...goqu.Expression) *goqu.SelectDataset {
	if len(expressions) == 1 {
		if sd, ok := expressions[0].(*goqu.SelectDataset); ok {
			return sd
		}
	}
	base := goqu.Select("*").From(tableName)
	if expressions == nil {
		return base
	}
	return base.Where(expressions...)
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
