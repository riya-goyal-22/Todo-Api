package repository

import (
	"database/sql"
	"errors"
	"todo/logger"
)

type UserRepository struct {
	client *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Validate(username, password string) bool {
	var name string
	row := r.client.QueryRow("SELECT username FROM users WHERE username = ? AND password=?", username, password)
	err := row.Scan(&name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("Error no rows found for user " + username)
			return false
		}
		logger.Log.Error("User " + username + " does not exist")
		return false
	}
	return true
}
