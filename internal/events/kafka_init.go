package events

import (
	"fmt"
	"net"
	"strconv"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// InitKafkaTopics ensures all necessary Kafka topics exist.
func InitKafkaTopics(cfg *config.KafkaConfig, logger *logrus.Logger) error {
	// Connect to any broker
	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()

	// Find controller node (handles topic operations)
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to controller at %s: %w", controllerAddr, err)
	}
	defer controllerConn.Close()

	// List of topics to create
	topics := []kafka.TopicConfig{
		{
			Topic:             cfg.Topics.UserEvents,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
		// Add more topics here
		// {Topic: cfg.Topics.SomeOtherTopic, NumPartitions: 3, ReplicationFactor: 1},
	}

	logger.Infof("Creating Kafka topics if they don't exist: %v", getTopicNames(topics))

	// Create topics
	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	logger.Info("Kafka topics initialized successfully")
	return nil
}

func getTopicNames(topics []kafka.TopicConfig) []string {
	names := make([]string, len(topics))
	for i, t := range topics {
		names[i] = t.Topic
	}
	return names
}
