package services

import (
	"database/sql"
	"github.com/google/uuid"
	"streakr-backend/pkg/utils"
)

type NewUser struct {
	Email      string `json:"email" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
	Password   string `json:"password" binding:"required"`
	PrivateKey string
	PublicKey  string
	Token      string
}

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email,omitempty"`
	APIKey         string `json:"api_key"`
	HashedPassword string `json:"hashed_password,omitempty"`
}

type RegistrationData struct {
	User *NewUser `json:"user" binding:"required"`
}

type UserLookup struct {
	Email string `json:"email"`
}

func (user *NewUser) Insert(db *sql.DB) (string, error) {
	u := &utils.Hash{}
	hash, err := u.Generate(user.Password)

	if err != nil {
		return "", err
	}

	var lastInsertId string
	err = db.QueryRow(
		"INSERT INTO users(id, email, api_key, password) VALUES($1, $2, $3, $4) RETURNING id;",
		uuid.New().String(),
		user.Email,
		user.APIKey,
		hash,
	).Scan(&lastInsertId)

	if err != nil {
		return "", err
	}

	return lastInsertId, nil
}

func (userLookup *UserLookup) GetByEmail(db *sql.DB) (User, error) {
	var user User

	err := db.QueryRow(
		"SELECT id, email, api_key FROM users WHERE email = $1",
		userLookup.Email,
	).Scan(&user.Id, &user.Email, &user.APIKey)

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
		err = rows.Scan(&user.Id, &user.Email)
		if err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}