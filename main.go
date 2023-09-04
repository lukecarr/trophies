package main

import (
	"github.com/lukecarr/trophies/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
