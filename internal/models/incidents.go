package models

import (
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

func ValidateIncident(i *Incident) error {
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
	return nil
}
