package routes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/internal/env"
	"net/http"
	"strconv"
)

func Metadata(e *env.Env, h *httprouter.Router) {
	h.GET("/api/metadata", getMetadata(e))
}

type MetadataResponse struct {
	Name            string `json:"name"`
	BackgroundImage string `json:"backgroundImage"`
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

		cacheKey := fmt.Sprintf("game.%d.background", game.ID)
		if x, found := e.Cache.Get(cacheKey); found {
			background := x.(string)
			if background != "missing" {
				resp.BackgroundImage = x.(string)
			}
		} else {
			metadata, err := e.Services.Metadata.SearchGame(game.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if metadata != nil {
				resp.BackgroundImage = metadata.BackgroundImageURL
				e.Cache.Set(cacheKey, resp.BackgroundImage, 0)
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
