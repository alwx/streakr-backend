package services

import (
	"database/sql"
	"streakr-backend/pkg/utils"

	"github.com/google/uuid"
)

type NewUser struct {
	Email        string `json:"email" binding:"required"`
	APIKey       string `json:"api_key" binding:"required"`
	Password     string `json:"password" binding:"required"`
	PublicKey    string
	PrivateKey   string
	Token        string
	DisplayName  string
	UserPersonId int64
}

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email,omitempty"`
	APIKey         string `json:"api_key,omitempty"`
	HashedPassword string `json:"hashed_password,omitempty"`
	PublicKey      string `json:"public_key,omitempty"`
	PrivateKey     string `json:"private_key,omitempty"`
	Token          string `json:"user_token,omitempty"`
	DisplayName    string `json:"display_name,omitempty"`
	UserPersonId   int    `json:"user_person_id,omitempty"`

	Campaigns []CampaignInfo `json:"campaigns,omitempty"`
}

type CampaignInfo struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Streak   int    `json:"streak"`
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
		"INSERT INTO users(id, email, api_key, public_key, private_key, user_token, password, display_name, user_person_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;",
		uuid.New().String(),
		user.Email,
		user.APIKey,
		user.PublicKey,
		user.PrivateKey,
		user.Token,
		hash,
		user.DisplayName,
		user.UserPersonId,
	).Scan(&lastInsertId)

	if err != nil {
		return "", err
	}

	return lastInsertId, nil
}

func (user *User) UpdateToken(db *sql.DB) error {
	return db.QueryRow(
		"UPDATE users SET user_token = $1 WHERE id = $2;",
		user.Token,
		user.Id,
	).Scan()
}

func (userLookup *UserLookup) GetByEmail(db *sql.DB) (User, error) {
	var user User

	err := db.QueryRow(
		"SELECT id, email, api_key, public_key, private_key, user_token, display_name, user_person_id FROM users WHERE email = $1",
		userLookup.Email,
	).Scan(&user.Id, &user.Email, &user.APIKey, &user.PublicKey, &user.PrivateKey, &user.Token, &user.DisplayName, &user.UserPersonId)

	if err != nil {
		return User{}, err
	}

	rows, err := db.Query(
		"SELECT c.name, c.badge_image_url, cu.streak_length FROM users as u LEFT JOIN campaign_user as cu ON cu.userId = u.id LEFT JOIN campaigns as c ON cu.campaignId = c.id  WHERE u.id = $1",
		user.Id,
	)

	if err != nil {
		return User{}, err
	}

	var campaigns []CampaignInfo
	defer rows.Close()
	for rows.Next() {
		var info CampaignInfo
		err = rows.Scan(&info.Name, &info.ImageUrl, &info.Streak)
		if err == nil {
			campaigns = append(campaigns, info)
		}
	}

	user.Campaigns = campaigns

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
