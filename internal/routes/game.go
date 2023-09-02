package routes

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/internal/env"
	"net/http"
)

func Game(e *env.Env, h *httprouter.Router) {
	h.GET("/api/games", getAllGames(e))
	h.GET("/api/games/counts", getAllGameTrophyCounts(e))
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
