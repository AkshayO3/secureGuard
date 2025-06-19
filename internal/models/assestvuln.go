package models

import (
	"database/sql"
	"errors"
	"time"
)

type AssetVuln struct {
	Id              string
	AssetId         string
	VulnerabilityId string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func ValidateAssetExists(db *sql.DB, a string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM assets WHERE id=$1)`
	err := db.QueryRow(query, a).Scan(&exists)
	if err != nil {
		return false
	}
	return true
}

func ValidateVulnerabilityExists(db *sql.DB, a string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM vulnerabilities WHERE id=$1)`
	err := db.QueryRow(query, a).Scan(&exists)
	if err != nil {
		return false
	}
	return true
}

func ValidateAssetVulnerability(db *sql.DB, av *AssetVuln) error {
	if !ValidateAssetExists(db, av.AssetId) {
		return errors.New("asset must exist in the database")
	}
	if !ValidateVulnerabilityExists(db, av.VulnerabilityId) {
		return errors.New("vulnerability must exist in the database")
	}
	if av.Status == "" {
		return errors.New("status cannot be empty")
	}
	return nil
}
