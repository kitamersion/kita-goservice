package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kitamersion/go-goservice/internal/api/handlers"
	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/database"
	"github.com/kitamersion/go-goservice/internal/domain/repositories"
	"github.com/kitamersion/go-goservice/internal/domain/services"
	"github.com/kitamersion/go-goservice/internal/events"
	"github.com/kitamersion/go-goservice/internal/events/producer"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Kafka topics (runs every startup, safe if already exists)
	if err := events.InitKafkaTopics(&cfg.Kafka, logger); err != nil {
		logger.WithError(err).Fatal("Failed to initialize Kafka topics")
	}

	// Initialize producer
	eventProducer := producer.NewProducer(&cfg.Kafka, logger)
	defer eventProducer.Close()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, eventProducer)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Setup Gin router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
		api.GET("/users/:id", userHandler.GetUser)
	}

	// Start server
	logger.Info("Starting API server on port ", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
