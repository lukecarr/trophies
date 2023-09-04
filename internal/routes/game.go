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

		if !game.BackgroundImageURL.Valid && !game.MetacriticScore.Valid && !game.ReleaseDate.Valid {
			platforms, excludePlatforms := platformToRawgPlatform(game.Platform)
			metadata, err := e.Services.Metadata.SearchGame(game.Name, platforms, excludePlatforms)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Update the SQL game table with new metadata
			updateQuery := "UPDATE game SET backgroundImageURL = $1, releaseDate = $2, metacriticScore = $3 WHERE id = $4"
			_, err = e.DB.Sql.Exec(updateQuery, metadata.BackgroundImageURL, metadata.ReleaseDate, metadata.MetacriticScore, id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			game.BackgroundImageURL = models.JsonNullString{NullString: sql.NullString{String: metadata.BackgroundImageURL, Valid: true}}
			game.ReleaseDate = models.JsonNullString{NullString: sql.NullString{String: metadata.ReleaseDate, Valid: true}}
			game.MetacriticScore = models.JsonNullInt64{NullInt64: sql.NullInt64{Int64: int64(metadata.MetacriticScore), Valid: true}}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(game); err != nil {
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
