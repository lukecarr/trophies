package models

import "database/sql"

type Game struct {
	ID             uint           `json:"id"`
	PsnID          string         `json:"psnID"`
	PsnServiceName string         `json:"psnServiceName"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	IconURL        string         `json:"iconURL"`
	Platform       string         `json:"platform"`
}
