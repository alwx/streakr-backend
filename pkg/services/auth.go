package services

import (
	"database/sql"
	"github.com/spf13/viper"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"github.com/go-redis/redis"
	"streakr-backend/pkg/utils"
)

func GetAuthMiddleware(db *sql.DB, redis *redis.Client) (*jwt.GinJWTMiddleware, error) {
	companyName, err := redis.Get("config.company_name").Result()
	if err != nil {
		companyName = ""
	}

	middleware := &jwt.GinJWTMiddleware{
		Realm:         companyName,
		Key:           []byte(companyName), //TODO(alwx): should be fixed
		Timeout:       time.Minute * time.Duration(viper.GetInt("web.auth.minutes_timeout")),
		MaxRefresh:    time.Hour * time.Duration(viper.GetInt("web.auth.hours_max_refresh")),
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Authenticator: func(email string, password string, c *gin.Context) (interface{}, bool) {
			login := Login{email, password}
			user, err := login.TryToLogin(db)

			if err != nil {
				return nil, false
			}
			return user, true
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"message": message})
		},
	}

	return middleware, nil
}

func ExtractJWTUser(c *gin.Context, db *sql.DB) (User, error) {
	claims := jwt.ExtractClaims(c)
	userEmail := claims["id"].(string)
	return (&UserLookup{Email: userEmail}).GetByEmail(db)
}


type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (login *Login) TryToLogin(db *sql.DB) (User, error) {
	var user User
	err := db.QueryRow(
		"SELECT username, email, password FROM users WHERE email = $1",
		login.Email,
	).Scan(&user.Username, &user.Email, &user.HashedPassword)

	u := &utils.Hash{}
	err = u.Compare(user.HashedPassword, login.Password)

	if err != nil {
		return User{}, err
	}
	return user, nil
}
