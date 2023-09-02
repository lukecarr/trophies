package env

import (
	"github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/services"
)

type Services struct {
	Game services.GameService
}

type Env struct {
	Services *Services
}

func New() *Env {
	return &Env{}
}

func NewSQlServices(db *db.DB) *Services {
	return &Services{
		Game: services.GameServiceSql{Sqlx: db.Sqlx},
	}
}
