package db

import (
	"github.com/jmoiron/sqlx"
)

const DIALECT = "sqlite3"
const MEMORY_DSN = ":memory:?cache=shared"

type DB struct {
	Sqlx *sqlx.DB
}

func New(dsn string) (*DB, error) {
	conn, err := sqlx.Connect(DIALECT, dsn)

	if err != nil {
		return nil, err
	}

	return &DB{
		Sqlx: conn,
	}, nil
}
