package models

import (
	"database/sql"
	"errors"
	"time"
)

type Incident struct {
	Id          string
	Title       string
	Description string
	Category    string
	Severity    string
	Status      string
	ReportedBy  string
	AssignedTo  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ValidateIncident(db *sql.DB, i *Incident) error {
	if i.Title == "" {
		return errors.New("title cannot be empty")
	}
	if i.Description == "" {
		return errors.New("description cannot be empty")
	}
	if i.Category == "" {
		return errors.New("category cannot be empty")
	}
	ValidSeverity := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
	if !ValidSeverity[i.Severity] {
		return errors.New("not a valid value for incident severity")
	}
	ValidStatus := map[string]bool{"detected": true, "investigating": true, "contained": true, "resolved": true}
	if !ValidStatus[i.Status] {
		return errors.New("not a valid value for incident status")
	}
	if !ValidateUserExists(db, i.ReportedBy) {
		return errors.New("reporting user was not found in records")
	}
	if i.AssignedTo != "" && !ValidateUserExists(db, i.AssignedTo) {
		return errors.New("assigned user was not found")
	}
	return nil
}
