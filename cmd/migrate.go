package cmd

import (
	"database/sql"
	"github.com/lukecarr/trophies/internal/db"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func MakeMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Performs database migrations (run after installation/upgrade)",
		Run: func(cmd *cobra.Command, args []string) {
			dsn := db.MemoryDsn
			if val, ok := os.LookupEnv("DSN"); ok {
				dsn = val
			}

			if dsn == "" {
				log.Fatalln("Please supply a path to a database file (SQLite) as an environment variable (DSN)")
			}

			conn, err := sql.Open(db.DIALECT, dsn)

			if err != nil {
				log.Fatalln("Failed to open SQLite connection", err)
			}

			n, err := db.Migrate(conn)

			if err != nil {
				log.Fatalln("Failed to perform SQLite migrations", err)
			} else if n == 0 {
				log.Println("The provided SQLite database is already up to date!")
			} else {
				log.Printf("Applied %d migration(s) successfully!\n", n)
			}
		},
	}
}

var migrateCmd = MakeMigrateCmd()

func init() {
	rootCmd.AddCommand(migrateCmd)
}
