package main

import (
	"context"
	"eci-backend/internal/auth/handler"
	"eci-backend/internal/auth/repository"
	"eci-backend/internal/auth/service"
	"eci-backend/internal/platform/config"
	"eci-backend/internal/platform/database"
	"eci-backend/internal/platform/middleware"
	"eci-backend/internal/platform/security"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI is required")
	}

	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ECI API is running",
		})
	})

	userRepository := repository.NewUserRepository(
		mongoClient.Database(cfg.MongoDatabase),
	)

	turnstileService := security.NewTurnstileService(
		cfg.TurnstileSecretKey,
	)

	db := mongoClient.Database(cfg.MongoDatabase)

	sessionRepository := repository.NewSessionRepository(db)

	authService := service.NewAuthService(userRepository,
		turnstileService, sessionRepository)
	authHandler := handler.NewAuthHandler(authService)

	app.POST("/api/auth/login", authHandler.Login)

	protected := app.Group("/api")
	protected.Use(middleware.Auth(sessionRepository))

	protected.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"userID": c.GetString("userID"),
		})
	})

	app.POST("/api/auth/logout", authHandler.Logout)

	protected.GET("/auth/me", authHandler.Me)

	app.Use(middleware.CORS())

	if err := app.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
