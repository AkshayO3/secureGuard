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

func (m *UserModel) ListUsers() ([]*models.User, error) {
	rows, err := m.DB.Query(`SELECT id, username, email, role, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (m *UserModel) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := m.DB.QueryRow(
		`SELECT id, username, email, role, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&user.Id, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) UpdateUserByID(id, username, email, role string) error {
	user, err := m.GetUserByID(id)
	if err != nil {
		return err
	}
	if username == "" {
		username = user.Username
	}
	if email == "" {
		email = user.Email
	}
	if role == "" {
		role = user.Role
	}
	_, err = m.DB.Exec(
		`UPDATE users SET username = $1, email = $2, role = $3, updated_at = NOW() WHERE id = $4`,
		username, email, role, id,
	)
	return err
}

func (m *UserModel) DeleteUserByID(id string) error {
	result, err := m.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}
