package models

import (
	"database/sql"
	"encoding/json"
)

type Game struct {
	ID             uint           `json:"id"`
	PsnID          string         `json:"psnID"`
	PsnServiceName string         `json:"psnServiceName"`
	Name           string         `json:"name"`
	Description    JsonNullString `json:"description"`
	IconURL        string         `json:"iconURL"`
	Platform       string         `json:"platform"`
}

type GameMetadata struct {
	GameID             uint           `json:"gameID"`
	BackgroundImageURL JsonNullString `json:"backgroundImageURL"`
	MetacriticScore    JsonNullInt64  `json:"metacriticScore"`
	ReleaseDate        JsonNullString `json:"releaseDate"`
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}
