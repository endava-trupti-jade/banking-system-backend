package main

import (
	_ "banking-system-backend/docs"
	"banking-system-backend/internal/config"
	"banking-system-backend/internal/middlewares"
	"banking-system-backend/internal/routes"
	"banking-system-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Banking System API
// @version 1.0
// @description Simple Banking System POC using Golang, Gin, MongoDB
// @termsOfService http://swagger.io/terms

// @host localhost:8080
// @BasePath /api
// @schemes http

// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log.Println("Hello banking poc backend....")

	// Load configuration
	config.Load()

	// Initialize Logger
	logger.InitLogger(config.AppConfig.AppEnv, config.AppConfig.AppName)
	defer logger.Sync()

	// Connect to Mongo
	mongoURI := "mongodb://" + config.AppConfig.DBHost + ":" + config.AppConfig.DBPort
	config.ConnectMongo(mongoURI, config.AppConfig.DBName)
	config.CreateIndexes()

	// Connect to Redis
	redisURI := config.AppConfig.RedisHost + ":" + config.AppConfig.RedisPort
	config.ConnectRedis(redisURI)

	//router := gin.Default()
	router := gin.New()
	// router.Use(gin.Recovery())
	// router.Use(middlewares.CORSMiddleware())
	// router.Use(middlewares.RateLimitMiddleware())
	router.Use(
		middlewares.RequestIDMiddleware(),
		middlewares.LoggerMiddleware(),
		middlewares.RecoveryMiddleware(),
		//middlewares.CORSMiddleware(),
		//middlewares.RateLimitMiddleware(),

	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(router)

	logger.Log.Info("application_started")

	router.Run(":" + config.AppConfig.AppPort)
}
