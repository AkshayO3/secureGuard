package models

import (
	"database/sql"
	"errors"
	"time"
)

type IncidentAsset struct {
	Id         string
	IncidentId string
	AssetId    string
	CreatedAt  time.Time
}

func ValidateIncidentExists(db *sql.DB, a string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM incidents WHERE id=$1)`
	err := db.QueryRow(query, a).Scan(&exists)
	if err != nil {
		return false
	}
	return true
}

func ValidateIncidentAsset(db *sql.DB, ia *IncidentAsset) error {
	if !ValidateAssetExists(db, ia.AssetId) {
		return errors.New("incident must exist in the db")
	}
	if !ValidateIncidentExists(db, ia.IncidentId) {
		return errors.New("incident must exist in the db")
	}
	return nil
}
