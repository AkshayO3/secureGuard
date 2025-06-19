package data

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"secureGuard/internal/models"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Register(username, email, role, password string) (string, error) {
	var exists bool
	err := m.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("user already exists")
	}
	user := &models.User{
		Username:     username,
		Email:        email,
		Role:         role,
		PasswordHash: password,
	}
	if err := models.ValidateUser(m.DB, user); err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	var id string
	err = m.DB.QueryRow(
		"INSERT INTO users (username, email, role, password_hash) VALUES ($1, $2, $3, $4) RETURNING id",
		username, email, role, string(hashed),
	).Scan(&id)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"userId": id,
		"role":   role,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m *UserModel) Login(username, password string) (string, error) {
	var id, role, hashed string
	err := m.DB.QueryRow(`SELECT id,role,password_hash FROM users WHERE username=$1`, username).Scan(&id, &role, &hashed)
	if err != nil {
		return "", errors.New(err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return "", errors.New("wrong password")
	}
	claims := jwt.MapClaims{
		"userId": id,
		"role":   role,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m *UserModel) RefreshToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userId"] == nil {
		return "", errors.New("invalid claims")
	}
	newClaims := jwt.MapClaims{
		"userId": claims["userId"],
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	signed, err := newToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
