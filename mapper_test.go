package goopy

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type TestStruct struct {
	Id    int64  `db:"id"`
	Email string `db:"email"`
}

func (t TestStruct) Table() string {
	return "users"
}

func (t TestStruct) Pk() (string, any) {
	return "id", t.Id
}

func Conn() (*pgx.Conn, error) {
	dburl := os.Getenv("GOOPY_DATABASE")
	conn, err := pgx.Connect(context.TODO(), dburl)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(context.TODO(),
		`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, email VARCHAR(100) NOT NULL)`,
	)
	if err != nil {
		return conn, err
	}
	_, err = conn.Exec(context.TODO(),
		`DELETE FROM users WHERE Id IS NOT NULL `,
	)
	if err != nil {
		return conn, err
	}
	_, err = conn.Exec(context.TODO(),
		`ALTER SEQUENCE users_id_seq RESTART WITH 1`,
	)

	return conn, err
}

func TestQuerier_Select(t *testing.T) {
	conn, err := Conn()
	assert.Nil(t, err)
	_, err = conn.Exec(context.TODO(), "INSERT INTO \"users\" (\"email\") VALUES ('email1@email.com'),('email2@email.com') ")
	assert.Nil(t, err)
	qr := NewPGQuerier[TestStruct](conn)
	res, err := qr.Select(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, []TestStruct{{int64(1), "email1@email.com"}, {int64(2), "email2@email.com"}}, res)
	res, err = qr.Select(context.TODO(), goqu.C("id").Eq(1))
	assert.Nil(t, err)
	assert.Equal(t, []TestStruct{{int64(1), "email1@email.com"}}, res)
}

func TestQuerier_SelectOne(t *testing.T) {
	conn, err := Conn()
	assert.Nil(t, err)
	_, err = conn.Exec(context.TODO(), "INSERT INTO \"users\" (\"email\") VALUES ('email1@email.com'),('email2@email.com') ")
	assert.Nil(t, err)
	qr := NewPGQuerier[TestStruct](conn)
	res, err := qr.SelectOne(context.TODO(), goqu.C("id").Eq(1))
	assert.Nil(t, err)
	assert.Equal(t, TestStruct{int64(1), "email1@email.com"}, res)

}

func TestMapper_Insert(t *testing.T) {
	conn, err := Conn()
	assert.Nil(t, err)
	qr := NewManager[TestStruct](conn)
	instance := TestStruct{
		Id:    1,
		Email: "email1@email.com",
	}
	res, err := qr.Insert(context.TODO(), instance)
	assert.Nil(t, err)
	assert.Equal(t, instance, res)
}

func TestMapper_Update(t *testing.T) {
	conn, err := Conn()
	assert.Nil(t, err)
	qr := NewManager[TestStruct](conn)
	instance := TestStruct{
		Id:    1,
		Email: "email1@email.com",
	}
	res, err := qr.Insert(context.TODO(), instance)
	assert.Nil(t, err)
	assert.Equal(t, instance, res)
	instance.Email = "email2@email.com"
	res, err = qr.Update(context.TODO(), instance)
	assert.Nil(t, err)
	assert.Equal(t, instance, res)
}

func TestMapper_Delete(t *testing.T) {
	conn, err := Conn()
	assert.Nil(t, err)
	qr := NewManager[TestStruct](conn)
	instance := TestStruct{
		Id:    1,
		Email: "email1@email.com",
	}
	res, err := qr.Insert(context.TODO(), instance)
	assert.Nil(t, err)
	assert.Equal(t, instance, res)
	err = qr.Delete(context.TODO(), instance)
	assert.Nil(t, err)
	_, err = qr.SelectOne(context.TODO(), goqu.C("id").Eq(1))
	assert.True(t, err.(Err).IsNotFound())
}
