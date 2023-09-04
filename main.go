package main

import (
	"github.com/lukecarr/trophies/cmd"
	"github.com/lukecarr/trophies/internal/info"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	if _, ok := os.LookupEnv("DISABLE_NEW_VERSION_CHECK"); !ok && info.Version != "undefined" {
		checkForLatestVersion()
	}

	cmd.Execute()
}

func checkForLatestVersion() {
	latestVersion, releaseDate, err := info.GetLatestVersion()

	if err != nil {
		log.Println("WARNING: Failed to check for newer version from GitHub!")
	}

	// Compare the latest version to the current version
	if latestVersion != info.Version {
		log.Printf("WARNING: A newer version of Trophies.gg is available: %s was released on %s\n!", latestVersion, releaseDate)
	}
}
