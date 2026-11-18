package main

import (
	"context"
	"eci-backend/internal/auth/handler"
	"eci-backend/internal/auth/repository"
	"eci-backend/internal/auth/service"
	employeeHandler "eci-backend/internal/employee/handler"
	employeeRepository "eci-backend/internal/employee/repository"
	employeeService "eci-backend/internal/employee/service"
	historyHandler "eci-backend/internal/history/handler"
	historyRepository "eci-backend/internal/history/repository"
	historyService "eci-backend/internal/history/service"
	historyWorker "eci-backend/internal/history/worker"
	numberToWordsHandler "eci-backend/internal/numbertowords/handler"
	numberToWordsService "eci-backend/internal/numbertowords/service"
	"eci-backend/internal/platform/config"
	"eci-backend/internal/platform/database"
	"eci-backend/internal/platform/middleware"
	"eci-backend/internal/platform/security"
	regionHandler "eci-backend/internal/region/handler"
	regionRepository "eci-backend/internal/region/repository"
	regionService "eci-backend/internal/region/service"
	starHandler "eci-backend/internal/star/handler"
	starService "eci-backend/internal/star/service"
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

	db := mongoClient.Database(cfg.MongoDatabase)

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
			"https://eci-fe.vercel.app",
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

	outboxRepo := historyRepository.NewOutboxRepository(db)
	employeeRepo := employeeRepository.NewEmployeeRepository(db)
	employeeSvc := employeeService.NewEmployeeService(employeeRepo, outboxRepo)
	employeeHTTP := employeeHandler.NewEmployeeHandler(employeeSvc)

	protected.GET("/employees", employeeHTTP.FindAll)
	protected.GET("/employees/:id", employeeHTTP.FindByID)
	protected.POST("/employees", employeeHTTP.Create)
	protected.PUT("/employees/:id", employeeHTTP.Update)
	protected.DELETE("/employees/:id", employeeHTTP.Delete)

	esURL := os.Getenv("ELASTICSEARCH_URL")
	esIndex := os.Getenv("ELASTICSEARCH_INDEX")
	esAPIKey := os.Getenv("ELASTICSEARCH_API_KEY")

	if esURL != "" && esIndex != "" && esAPIKey != "" {
		esWorker := historyWorker.NewElasticsearchWorker(
			outboxRepo,
			esURL,
			esIndex,
			esAPIKey,
		)

		workerCtx, cancelWorker := context.WithCancel(context.Background())
		defer cancelWorker()

		go esWorker.Run(workerCtx)
	}

	esRepo := historyRepository.NewElasticsearchRepository(
		esURL,
		esIndex,
		esAPIKey,
	)

	historySvc := historyService.NewHistoryService(esRepo)
	historyHTTP := historyHandler.NewHistoryHandler(historySvc)

	protected.GET("/logs", historyHTTP.FindAll)

	numberToWordsSvc := numberToWordsService.NewNumberToWordsService(outboxRepo)
	numberToWordsHTTP := numberToWordsHandler.NewNumberToWordsHandler(numberToWordsSvc)

	protected.POST("/number-to-words", numberToWordsHTTP.Convert)

	starSvc := starService.NewStarService(outboxRepo)
	starHTTP := starHandler.NewStarHandler(starSvc)

	protected.POST("/star-pattern", starHTTP.Generate)

	regionRepo := regionRepository.NewRegionRepository(db)
	regionSvc := regionService.NewRegionService(regionRepo)
	regionHTTP := regionHandler.NewRegionHandler(regionSvc)

	protected.GET("/regions/provinces", regionHTTP.GetProvinces)
	protected.GET("/regions/provinces/:provinceId/cities", regionHTTP.GetCities)
	protected.GET("/regions/cities/:cityId/districts", regionHTTP.GetDistricts)

	if err := app.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
