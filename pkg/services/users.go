package services

import (
	"database/sql"
	"github.com/google/uuid"
	"streakr-backend/pkg/utils"
)

type NewUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email,omitempty"`
	HashedPassword string `json:"hashed_password,omitempty"`
}

type RegistrationData struct {
	User *NewUser `json:"user" binding:"required"`
}

func (user *NewUser) Insert(db *sql.DB) (string, error) {
	u := &utils.Hash{}
	hash, err := u.Generate(user.Password)

	if err != nil {
		return "", err
	}

	var lastInsertId string
	err = db.QueryRow(
		"INSERT INTO users(id, username, email, password) VALUES($1, $2, $3, $4) RETURNING id;",
		uuid.New().String(),
		user.Username,
		user.Email,
		hash,
	).Scan(&lastInsertId)

	if err != nil {
		return "", err
	}

	return lastInsertId, nil
}


type UserLookup struct {
	Email string `json:"email"`
}

func (userLookup *UserLookup) GetByEmail(db *sql.DB) (User, error) {
	var user User

	err := db.QueryRow(
		"SELECT id, username FROM users WHERE email = $1",
		userLookup.Email,
	).Scan(&user.Id, &user.Username)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUsers(db *sql.DB) ([]User, error) {
	var users []User

	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username)
		if err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}