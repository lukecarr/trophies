package cmd

import (
	sql "github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/services"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sync"
)

func MakeFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Fetches games, trophy lists, and trophy completion data from PSN.",
		Run: func(cmd *cobra.Command, args []string) {
			dsn, ok := os.LookupEnv("DSN")
			if !ok {
				log.Fatalln("You must set the DSN environment variable to a valid SQLite connection string")
			}

			conn, err := sql.New(dsn)

			if err != nil {
				log.Fatalln("Failed to open SQLite connection", err)
			}

			npsso, ok := os.LookupEnv("NPSSO")
			if !ok {
				log.Fatalln("You must set the NPSSO environment variable")
			}

			psn := services.NewPsnClient(npsso)
			titles, _ := psn.GetTitles()

			tx, err := conn.Sql.Begin()
			if err != nil {
				log.Fatalln("Failed to begin transaction for loading games:", err)
			}

			query := `
				INSERT INTO game (name, iconURL, description, psnID, psnServiceName, platform)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT(psnID) DO UPDATE SET
				name = excluded.name,
				iconURL = excluded.iconURL,
				description = excluded.description,
				psnServiceName = excluded.psnServiceName,
				platform = excluded.platform
			`

			for _, title := range titles {
				_, err = tx.Exec(query, title.Name, title.IconURL, title.Description, title.ID, title.ServiceName, title.Platform)

				if err != nil {
					_ = tx.Rollback()
					log.Fatalln("Failed to load games:", err)
				}
			}

			allTrophyGroups := getTrophyGroups(titles, psn)

			trophyGroupQuery := `
				INSERT INTO trophyGroup (name, iconURL, gameID, psnID)
				VALUES ($1, $2, (SELECT id FROM game WHERE psnID = $3), $4)
				ON CONFLICT(gameID, psnID) DO UPDATE SET
				name = excluded.name,
				iconURL = excluded.iconURL
			`

			for _, trophyGroup := range allTrophyGroups {
				_, err = tx.Exec(trophyGroupQuery, trophyGroup.Name, trophyGroup.IconURL, trophyGroup.GameID, trophyGroup.ID)
				if err != nil {
					_ = tx.Rollback()
					log.Fatalln("Failed to load trophy groups:", err)
				}
			}

			allTrophies := getTrophies(titles, psn)

			trophyQuery := `
				INSERT INTO trophy (name, description, rarity, hidden, iconURL, gameID, psnID, trophyGroupID)
				VALUES ($1, $2, $3, $4, $5, (SELECT id FROM game WHERE psnID = $6), $7, (SELECT id FROM trophyGroup WHERE psnID = $8))
				ON CONFLICT(gameID, psnID) DO UPDATE SET
				name = excluded.name,
				description = excluded.description,
				rarity = excluded.rarity,
				hidden = excluded.hidden,
				iconURL = excluded.iconURL
			`

			for _, trophy := range allTrophies {
				_, err = tx.Exec(trophyQuery, trophy.Name, trophy.Description, trophy.Rarity, trophy.Hidden, trophy.IconURL, trophy.GameID, trophy.ID, trophy.TrophyGroupID)
				if err != nil {
					_ = tx.Rollback()
					log.Fatalln("Failed to load trophies:", err)
				}
			}

			_ = tx.Commit()

			log.Printf("Loaded %d games, %d trophy groups, and %d trophies from PSN!\n", len(titles), len(allTrophyGroups), len(allTrophies))
		},
	}
}

type TrophyGroupRow struct {
	ID      string `db:"psnID"`
	Name    string `db:"name"`
	GameID  string `db:"psnGameID"`
	IconURL string `db:"iconURL"`
}

type TrophyRow struct {
	ID            uint   `db:"psnID"`
	Name          string `db:"name"`
	Description   string `db:"description"`
	IconURL       string `db:"iconURL"`
	Rarity        string `db:"rarity"`
	Hidden        bool   `db:"hidden"`
	TrophyGroupID string `db:"trophyGroupID"`
	GameID        string `db:"psnGameID"`
}

func getTrophyGroups(titles []services.Title, psn services.PsnService) []TrophyGroupRow {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(titles))

	allTrophyGroupRows := make([]TrophyGroupRow, 0)
	var trophyGroupMutex sync.Mutex

	for _, title := range titles {
		wg.Add(1)
		go func(title services.Title) {
			defer wg.Done()

			trophyGroups, err := psn.GetTrophyGroups(title.ID, title.ServiceName)
			if err != nil {
				errorChannel <- err
				return
			}

			rows := make([]TrophyGroupRow, len(trophyGroups))
			for i, trophyGroup := range trophyGroups {
				rows[i] = TrophyGroupRow{
					ID:      trophyGroup.ID,
					Name:    trophyGroup.Name,
					GameID:  title.ID,
					IconURL: trophyGroup.IconURL,
				}
			}

			trophyGroupMutex.Lock()
			allTrophyGroupRows = append(allTrophyGroupRows, rows...)
			trophyGroupMutex.Unlock()
		}(title)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		log.Println("Error:", err)
	}

	return allTrophyGroupRows
}

func getTrophies(titles []services.Title, psn services.PsnService) []TrophyRow {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(titles))

	allTrophyRows := make([]TrophyRow, 0)
	var trophyMutex sync.Mutex

	for _, title := range titles {
		wg.Add(1)
		go func(title services.Title) {
			defer wg.Done()

			trophies, err := psn.GetTrophies(title.ID, title.ServiceName)
			if err != nil {
				errorChannel <- err
				return
			}

			rows := make([]TrophyRow, len(trophies))
			for i, trophy := range trophies {
				rows[i] = TrophyRow{
					ID:            trophy.ID,
					Name:          trophy.Name,
					Description:   trophy.Description,
					Hidden:        trophy.Hidden,
					Rarity:        trophy.Rarity,
					GameID:        title.ID,
					TrophyGroupID: trophy.GroupID,
					IconURL:       trophy.IconURL,
				}
			}

			trophyMutex.Lock()
			allTrophyRows = append(allTrophyRows, rows...)
			trophyMutex.Unlock()
		}(title)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		log.Println("Error:", err)
	}

	return allTrophyRows
}

var fetchCmd = MakeFetchCmd()

func init() {
	rootCmd.AddCommand(fetchCmd)
}
