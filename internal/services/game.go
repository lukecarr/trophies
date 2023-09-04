package services

import (
	"database/sql"
	"github.com/lukecarr/trophies/internal/models"
)

type GameTrophyCounts struct {
	ID     int    `json:"id"`
	Rarity string `json:"rarity"`
	Count  int    `json:"count"`
}

type GameService interface {
	Get(id int) (*models.Game, error)
	GetAll() ([]*models.Game, error)
	GetCounts() ([]*GameTrophyCounts, error)
}

type GameServiceSql struct {
	Sql *sql.DB
}

const GetQuery = `
	SELECT
		id,
		psnID,
		name,
		description,
		iconURL,
		platform,
		psnServiceName
	FROM game
	WHERE id = $1
`

func (s GameServiceSql) Get(id int) (*models.Game, error) {
	game := &models.Game{}
	err := s.Sql.QueryRow(GetQuery, id).Scan(&game.ID, &game.PsnID, &game.Name, &game.Description, &game.IconURL, &game.Platform, &game.PsnServiceName)

	return game, err
}

const GetAllQuery = `
	SELECT
		id,
		psnID,
		name,
		description,
		iconURL,
		platform,
		psnServiceName
	FROM game
`

func (s GameServiceSql) GetAll() ([]*models.Game, error) {
	games := make([]*models.Game, 0)

	rows, err := s.Sql.Query(GetAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		game := &models.Game{}
		err := rows.Scan(&game.ID, &game.PsnID, &game.Name, &game.Description, &game.IconURL, &game.Platform, &game.PsnServiceName)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

const GetCountsQuery = `
	SELECT
		game.id,
		trophy.rarity,
		COUNT(trophy.id) AS count
	FROM game
	LEFT JOIN trophy ON trophy.gameID = game.id
	GROUP BY game.id, trophy.rarity
`

func (s GameServiceSql) GetCounts() ([]*GameTrophyCounts, error) {
	counts := make([]*GameTrophyCounts, 0)

	rows, err := s.Sql.Query(GetCountsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		count := &GameTrophyCounts{}
		err := rows.Scan(&count.ID, &count.Rarity, &count.Count)
		if err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}

	return counts, nil
}
