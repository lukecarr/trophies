package models

import "database/sql"

type Game struct {
	ID             uint           `json:"id" db:"id"`
	PsnID          string         `json:"psnID" db:"psnID"`
	PsnServiceName string         `json:"psnServiceName" db:"psnServiceName"`
	Name           string         `json:"name" db:"name"`
	Description    sql.NullString `json:"description" db:"description"`
	IconURL        string         `json:"iconURL" db:"iconURL"`
	Platform       string         `json:"platform" db:"platform"`
}
