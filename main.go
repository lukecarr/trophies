package main

import (
	"database/sql"
	"github.com/lukecarr/trophies/cmd"
	"github.com/lukecarr/trophies/internal/db"
	"modernc.org/sqlite"
)

func main() {
	sql.Register(db.DIALECT, &sqlite.Driver{})

	cmd.Execute()
}
