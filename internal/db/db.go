package db

import (
	"database/sql"
)

const DIALECT = "sqlite3"
const MemoryDsn = ":memory:?cache=shared"

type DB struct {
	Sql *sql.DB
}

func New(dsn string) (*DB, error) {
	conn, err := sql.Open(DIALECT, dsn)

	if err != nil {
		return nil, err
	}

	return &DB{
		Sql: conn,
	}, nil
}
