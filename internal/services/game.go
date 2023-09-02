package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/lukecarr/trophies/internal/models"
)

type GameTrophyCounts struct {
	ID     int    `json:"id" db:"id"`
	Rarity string `json:"rarity" db:"rarity"`
	Count  int    `json:"count" db:"count"`
}

type GameService interface {
	Get(id int) (*models.Game, error)
	GetAll() ([]*models.Game, error)
	GetCounts() ([]*GameTrophyCounts, error)
}

type GameServiceSql struct {
	Sqlx *sqlx.DB
}

func (s GameServiceSql) Get(id int) (*models.Game, error) {
	game := &models.Game{}
	err := s.Sqlx.Get(game, "SELECT * FROM game WHERE id = $1", id)
	return game, err
}

func (s GameServiceSql) GetAll() ([]*models.Game, error) {
	games := make([]*models.Game, 0)
	err := s.Sqlx.Select(&games, "SELECT * FROM game")
	return games, err
}

func (s GameServiceSql) GetCounts() ([]*GameTrophyCounts, error) {
	counts := make([]*GameTrophyCounts, 0)
	err := s.Sqlx.Select(&counts, `
		SELECT
			game.id,
			trophy.rarity,
			COUNT(trophy.id) AS count
		FROM game
		LEFT JOIN trophy ON trophy.gameID = game.id
		GROUP BY game.id, trophy.rarity
	`)
	return counts, err
}
