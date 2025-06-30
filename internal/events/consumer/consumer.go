package consumer

import (
	"context"
	"time"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// EventHandler handles proto messages with their raw JSON payload and headers
type EventHandler func(ctx context.Context, eventType string, headers map[string]string, payload []byte) error

type Consumer struct {
	reader   *kafka.Reader
	logger   *logrus.Logger
	handlers map[string]EventHandler // eventType -> handler mapping
}

func NewConsumer(cfg *config.KafkaConfig, logger *logrus.Logger) *Consumer {
	readerConfig := kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topics.UserEvents,
		GroupID:        cfg.ConsumerGroups.UserConsumer,
		StartOffset:    kafka.FirstOffset,
		MinBytes:       1,           // 1B
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
	}
	logger.Infof("Creating Kafka reader with config: Brokers=%v Topic=%s GroupID=%s", readerConfig.Brokers, readerConfig.Topic, readerConfig.GroupID)

	reader := kafka.NewReader(readerConfig)

	return &Consumer{
		reader:   reader,
		logger:   logger,
		handlers: make(map[string]EventHandler),
	}
}

func (c *Consumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka consumer loop")
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Context cancelled, shutting down consumer")
			return ctx.Err()
		default:
			c.logger.Debug("Waiting to read message...")
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.WithError(err).Error("Failed to read message")
				time.Sleep(time.Second) // backoff on error
				continue
			}

			c.logger.Infof("Received message at topic %s partition %d offset %d", message.Topic, message.Partition, message.Offset)
			c.logger.Debugf("Message key: %s, value: %s", string(message.Key), string(message.Value))

			// Extract headers from Kafka message
			headers := make(map[string]string)
			eventType := ""
			for _, header := range message.Headers {
				headers[header.Key] = string(header.Value)
				if header.Key == "event_type" {
					eventType = string(header.Value)
				}
			}

			c.logger.Infof("Event type: %s", eventType)

			if handler, exists := c.handlers[eventType]; exists {
				if err := handler(ctx, eventType, headers, message.Value); err != nil {
					c.logger.WithError(err).WithField("event_type", eventType).Error("Failed to handle event")
				} else {
					c.logger.WithField("event_type", eventType).Info("Event handled successfully")
				}
			} else {
				c.logger.WithField("event_type", eventType).Warn("No handler registered for event type")
			}
		}
	}
}

func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka reader")
	return c.reader.Close()
}
