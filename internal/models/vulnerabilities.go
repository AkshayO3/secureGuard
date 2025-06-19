package models

import (
	"errors"
	"time"
)

type Vulnerability struct {
	Id          string
	CveId       string
	Title       string
	Description string
	Severity    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ValidateVulnerability(v *Vulnerability) error {
	if v.Title == "" {
		return errors.New("title cannot be empty")
	}
	if v.Description == "" {
		return errors.New("description cannot be empty")
	}
	ValidSeverity := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
	if !ValidSeverity[v.Severity] {
		return errors.New("not a valid value for incident severity")
	}
	ValidStatus := map[string]bool{"open": true, "in_progress": true, "mitigated": true, "resolved": true}
	if !ValidStatus[v.Status] {
		return errors.New("not a valid value for incident status")
	}
	return nil
}
