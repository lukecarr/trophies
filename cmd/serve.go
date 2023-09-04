package cmd

import (
	"github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/server"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func MakeServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Starts Trophies.gg's web server",
		Run: func(cmd *cobra.Command, args []string) {
			npsso, ok := os.LookupEnv("NPSSO")
			if !ok {
				log.Fatalln("You must set the NPSSO environment variable")
			}

			dsn := db.MemoryDsn
			if val, ok := os.LookupEnv("DSN"); ok {
				dsn = val
			}

			addr := ":3000"
			if val, ok := os.LookupEnv("ADDR"); ok {
				addr = val
			}

			rawg := os.Getenv("RAWG_API_KEY")

			srv := server.New(dsn, npsso, rawg)
			srv.Listen(addr)
		},
	}
}

var serveCmd = MakeServeCmd()

func init() {
	rootCmd.AddCommand(serveCmd)
}
