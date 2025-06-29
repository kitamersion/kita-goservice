package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/events"
	"github.com/kitamersion/go-goservice/internal/events/consumer"
	"github.com/kitamersion/go-goservice/internal/events/consumer/handlers"
	"github.com/kitamersion/go-goservice/internal/events/types"
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

	logger.Info("User topic: ", cfg.Kafka.Topics.UserEvents)

	// Initialize Kafka topics (runs every startup, safe if already exists)
	if err := events.InitKafkaTopics(&cfg.Kafka, logger); err != nil {
		logger.WithError(err).Fatal("Failed to initialize Kafka topics")
	}

	// Initialize consumer
	eventConsumer := consumer.NewConsumer(&cfg.Kafka, logger)
	defer eventConsumer.Close()

	// Initialize event handlers
	userHandlers := handlers.NewUserEventHandlers(logger)

	// Register event handlers with the correct signature
	eventConsumer.RegisterHandler(types.UserCreated, userHandlers.HandleUserCreated)
	eventConsumer.RegisterHandler(types.UserUpdated, userHandlers.HandleUserUpdated)
	eventConsumer.RegisterHandler(types.UserDeleted, userHandlers.HandleUserDeleted)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Shutting down consumer...")
		cancel()
	}()

	// Start consuming
	logger.Info("Starting event consumer...")
	if err := eventConsumer.Start(ctx); err != nil {
		logger.WithError(err).Error("Consumer stopped with error")
	}
}
