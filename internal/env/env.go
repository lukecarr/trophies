package env

import (
	"github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/services"
	"github.com/patrickmn/go-cache"
	"time"
)

type Services struct {
	Game     services.GameService
	Metadata services.MetadataService
}

type Env struct {
	Services *Services
	Cache    *cache.Cache
}

func New() *Env {
	return &Env{
		Cache: cache.New(24*7*time.Hour, 1*time.Hour),
	}
}

func NewServices(db *db.DB, rawgApiKey string) *Services {
	return &Services{
		Game:     services.GameServiceSql{Sql: db.Sql},
		Metadata: &services.MetadataServiceRawg{ApiKey: rawgApiKey},
	}
}
