package models

import (
	"database/sql"
	"errors"
	"time"
)

type Audit struct {
	Id           string
	UserId       string
	Action       string
	ResourceType string
	ResourceId   string
	Details      string
	Ip           string
	CreatedAt    time.Time
}

func ValidateUserExists(db *sql.DB, a string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`
	err := db.QueryRow(query, a).Scan(&exists)
	if err != nil {
		return false
	}
	return true
}

func ValidateAudit(db *sql.DB, a *Audit) error {
	if !ValidateUserExists(db, a.UserId) {
		return errors.New("user must exist in the database")
	}
	if a.Action == "" {
		return errors.New("action cannot be empty")
	}
	if a.ResourceType == "" {
		return errors.New("resource type cannot be empty")
	}
	if a.Ip == "" {
		return errors.New("ip address cannot be empty")
	}
	return nil
}
