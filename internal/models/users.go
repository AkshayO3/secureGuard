package models

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
	Id           string
	Username     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func IsEmailUnique(db *sql.DB, email string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false
	}
	return !exists
}

func IsUsernameUnique(db *sql.DB, username string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false
	}
	return !exists
}

func ValidateUser(db *sql.DB, u *User) error {
	if u.Username == "" || !IsUsernameUnique(db, u.Username) {
		return errors.New("username must be unique and non-empty")
	}
	if u.Email == "" || !strings.HasSuffix(u.Email, "@gmail.com") || !IsEmailUnique(db, u.Email) {
		return errors.New("email must be unique, non-empty, and a valid Gmail address")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@gmail\.com$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("email must be a valid Gmail address")
	}
	if u.PasswordHash == "" {
		return errors.New("password hash must be non-empty")
	}
	validRoles := map[string]bool{"admin": true, "analyst": true, "viewer": true}
	if !validRoles[u.Role] {
		return errors.New("role must be one of admin, analyst, viewer")
	}
	return nil
}
