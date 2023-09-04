package env

import (
	"github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/services"
)

type Services struct {
	Game     services.GameService
	Metadata services.MetadataService
}

type Env struct {
	DB       *db.DB
	Services *Services
}

func New(db *db.DB) *Env {
	return &Env{
		DB: db,
	}
}

func NewServices(db *db.DB, rawgApiKey string) *Services {
	return &Services{
		Game:     services.GameServiceSql{Sql: db.Sql},
		Metadata: &services.MetadataServiceRawg{ApiKey: rawgApiKey},
	}
}
