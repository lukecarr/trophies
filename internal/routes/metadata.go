package routes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/internal/env"
	"net/http"
	"strconv"
	"strings"
)

func Metadata(e *env.Env, h *httprouter.Router) {
	h.GET("/api/metadata", getMetadata(e))
}

type MetadataResponse struct {
	Name            string `json:"name"`
	BackgroundImage string `json:"backgroundImage"`
	MetacriticScore int    `json:"metacritic"`
}

func getMetadata(e *env.Env) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "ID must be an integer", http.StatusBadRequest)
			return
		}

		game, err := e.Services.Game.Get(idInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := MetadataResponse{
			Name: game.Name,
		}

		cacheKey := fmt.Sprintf("game.%d.metadata", game.ID)
		if x, found := e.Cache.Get(cacheKey); found {
			err := json.Unmarshal(x.([]byte), &resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			platforms, excludePlatforms := platformToRawgPlatform(game.Platform)
			metadata, err := e.Services.Metadata.SearchGame(game.Name, platforms, excludePlatforms)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if metadata != nil {
				resp.BackgroundImage = metadata.BackgroundImageURL
				resp.MetacriticScore = metadata.MetacriticScore

				jsonString, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				e.Cache.Set(cacheKey, jsonString, 0)
			} else {
				e.Cache.Set(cacheKey, "missing", 0)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
