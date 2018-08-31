package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"github.com/go-redis/redis"
	"github.com/appleboy/gin-jwt"
	"streakr-backend/pkg/services"
	"time"
)

type Data struct {
	Database *sql.DB
	Redis *redis.Client
	AuthMiddleware *jwt.GinJWTMiddleware
	Router *gin.Engine
	SecureArea *gin.RouterGroup
}

func InitHttp(db *sql.DB, redis *redis.Client) {
	authMiddleware, err := services.GetAuthMiddleware(db, redis)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	data := Data{
		Database: db,
		Redis: redis,
		AuthMiddleware: authMiddleware,
		Router: router,
	}

	UserRouter(data)
	secureArea := router.Group("")
	secureArea.Use(authMiddleware.MiddlewareFunc())
	{
		data.SecureArea = secureArea
		// more routers come here
	}

	router.Run(fmt.Sprintf(":%s", viper.GetString("http.port")))
}
