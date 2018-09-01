package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"github.com/appleboy/gin-jwt"
	"streakr-backend/pkg/services"
	"time"
)

type Data struct {
	Database *sql.DB
	AuthMiddleware *jwt.GinJWTMiddleware
	Router *gin.Engine
	SecureArea *gin.RouterGroup
}

func InitHttp(db *sql.DB) {
	authMiddleware, err := services.GetAuthMiddleware(db)
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
