package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func MakeRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "trophies",
		Short: "Run and manage Trophies.gg, a self-hosted web application for tracking PlayStation trophies ",
	}
}

var rootCmd = MakeRootCmd()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			log.Panicln(err)
		}

		os.Exit(1)
	}
}
