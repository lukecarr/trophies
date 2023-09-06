package routes

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/internal/env"
	"github.com/lukecarr/trophies/internal/models"
	"net/http"
	"strconv"
	"strings"
)

func Game(e *env.Env, h *httprouter.Router) {
	h.GET("/api/games", getAllGames(e))
	h.GET("/api/games/:id", getGame(e))
	h.GET("/api/games/:id/metadata", getGameMetadata(e))
	h.GET("/api/gamesCounts", getAllGameTrophyCounts(e))
}

func getAllGames(e *env.Env) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		games, err := e.Services.Game.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(games); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func platformToRawgPlatform(platforms string) (string, string) {
	if platforms == "" {
		return "", ""
	}

	platformsMap := map[string]string{
		"PS3":    "16",
		"PSVITA": "19",
		"PS4":    "18",
		"PS5":    "187",
	}

	allPlatforms := strings.Split(platforms, ",")

	for _, platform := range allPlatforms {
		if rawgPlatform, ok := platformsMap[platform]; ok {
			// Build a string of all other platforms (excluding the matched one)
			otherPlatforms := make([]string, len(allPlatforms)-1)
			for _, x := range allPlatforms {
				if x != platform {
					otherPlatforms = append(otherPlatforms, x)
				}
			}

			return rawgPlatform, strings.Join(otherPlatforms, ",")
		}
	}

	return "", ""
}

func getGame(e *env.Env) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idString := p.ByName("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		game, err := e.Services.Game.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(game); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

const InsertMetadataQuery = `
	INSERT INTO gameMetadata ("gameID", "backgroundImageURL", "metacriticScore", "releaseDate") VALUES ($1, $2, $3, $4)
`

func getGameMetadata(e *env.Env) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idString := p.ByName("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		game, err := e.Services.Game.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metadata, err := e.Services.Game.GetMetadata(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if metadata == nil {
			platforms, excludePlatforms := platformToRawgPlatform(game.Platform)
			newMetadata, err := e.Services.Metadata.SearchGame(game.Name, platforms, excludePlatforms)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = e.DB.Sql.Exec(InsertMetadataQuery, id, newMetadata.BackgroundImageURL, newMetadata.MetacriticScore, newMetadata.ReleaseDate)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			metadata = &models.GameMetadata{}
			metadata.GameID = game.ID
			metadata.BackgroundImageURL = models.JsonNullString{NullString: sql.NullString{String: newMetadata.BackgroundImageURL, Valid: true}}
			metadata.MetacriticScore = models.JsonNullInt64{NullInt64: sql.NullInt64{Int64: int64(newMetadata.MetacriticScore), Valid: true}}
			metadata.ReleaseDate = models.JsonNullString{NullString: sql.NullString{String: newMetadata.ReleaseDate, Valid: true}}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(metadata); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getAllGameTrophyCounts(e *env.Env) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		games, err := e.Services.Game.GetCounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(games); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
