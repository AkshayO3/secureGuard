package models

import (
	"errors"
	"time"
)

type Asset struct {
	Id        string
	Name      string
	Type      string
	Ip        string
	Os        string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ValidateAsset(a *Asset) error {
	if a.Name == "" {
		return errors.New("name cannot be empty")
	}
	ValidType := map[string]bool{"server": true, "workstation": true, "network_device": true, "application": true}
	if !ValidType[a.Type] {
		return errors.New("not a valid value for incident severity")
	}
	ValidStatus := map[string]bool{"vulnerable": true, "patched": true, "not_applicable": true}
	if !ValidStatus[a.Status] {
		return errors.New("not a valid value for incident status")
	}
	return nil
}
