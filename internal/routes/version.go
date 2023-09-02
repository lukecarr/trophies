package routes

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/internal/info"
	"log"
	"net/http"
	"time"
)

type versionResponse struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
}

func Version() func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	res := versionResponse{}

	if info.Version != "undefined" {
		res.Version = info.Version
	}

	date := info.Date
	if date == "undefined" {
		date = time.Now().Format(time.RFC3339)
	}
	res.Date = date

	if info.Commit != "undefined" {
		res.Commit = info.Commit
	}

	return func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err.Error())
			_, _ = fmt.Fprint(w, "{}")
		}
	}
}
