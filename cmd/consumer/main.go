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

	// TODO: make this generic for additional consumers to get registered
	kafkaConsumer := consumer.NewConsumer(&cfg.Kafka, logger)
	defer kafkaConsumer.Close()

	// Initialize event handlers
	userHandlers := handlers.NewUserEventHandlers(logger)

	// Register handlers for each proto event type
	// The event type string matches what the producer sends in the "event_type" header
	kafkaConsumer.RegisterHandler("userpb.UserCreated", userHandlers.HandleUserCreated)
	kafkaConsumer.RegisterHandler("userpb.UserUpdated", userHandlers.HandleUserUpdated)
	kafkaConsumer.RegisterHandler("userpb.UserDeleted", userHandlers.HandleUserDeleted)

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
	if err := kafkaConsumer.Start(ctx); err != nil {
		logger.WithError(err).Error("Consumer stopped with error")
	}
}
