package goopy

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

type userStruct struct {
	Id    int    `db:"id"`
	Email string `db:"email"`
}

func (u userStruct) Table() string {
	return "users"
}

func (u userStruct) Pk() (string, any) {
	return "id", u.Id
}

func TestQuerier_toSelectQuery(t *testing.T) {
	builder := toSelectQuery("users", goqu.C("uuid").Eq(1))
	sql, args, _ := builder.Prepared(true).ToSQL()
	assert.Equal(t, "SELECT * FROM \"users\" WHERE (\"uuid\" = ?)", sql)
	assert.Equal(t, []any{int64(1)}, args)
}

func TestQuerier_toInsertQuery(t *testing.T) {
	userStruct := struct {
		Id    int64  `db:"id"`
		Email string `db:"email"`
	}{}
	user := userStruct
	user.Id = 1
	user.Email = "email@email.com"
	builder := toInsertQuery("users", user)
	sql, args, _ := builder.Prepared(true).ToSQL()
	assert.Equal(t, "INSERT INTO \"users\" (\"email\", \"id\") VALUES (?, ?) RETURNING *", sql)
	assert.Equal(t, []any{"email@email.com", int64(1)}, args)
}

func TestQuerier_toUpdateQuery(t *testing.T) {
	var user userStruct
	user.Id = 1
	user.Email = "email@email.com"
	builder := toUpdateQuery("users", user)
	sql, args, _ := builder.Prepared(true).ToSQL()
	assert.Equal(t, "UPDATE \"users\" SET \"email\"=?,\"id\"=? WHERE (\"id\" = ?) RETURNING *", sql)
	assert.Equal(t, []any{"email@email.com", int64(1), int64(1)}, args)
}

func TestQuerier_toDeleteQuery(t *testing.T) {
	var user userStruct
	user.Id = 1
	user.Email = "email@email.com"
	builder := toDeleteQuery("users", user)
	sql, args, _ := builder.Prepared(true).ToSQL()
	assert.Equal(t, "DELETE FROM \"users\" WHERE (\"id\" = ?)", sql)
	assert.Equal(t, []any{int64(1)}, args)
}
