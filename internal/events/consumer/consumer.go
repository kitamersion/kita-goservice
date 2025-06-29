package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/events/types"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type EventHandler func(ctx context.Context, event *types.BaseEvent) error

type Consumer struct {
	reader   *kafka.Reader
	logger   *logrus.Logger
	handlers map[types.EventType]EventHandler
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
		handlers: make(map[types.EventType]EventHandler),
	}
}

func (c *Consumer) RegisterHandler(eventType types.EventType, handler EventHandler) {
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

			var event types.BaseEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				c.logger.WithError(err).Error("Failed to unmarshal event")
				continue
			}

			c.logger.Infof("Event type: %s", event.Type)

			if handler, exists := c.handlers[event.Type]; exists {
				if err := handler(ctx, &event); err != nil {
					c.logger.WithError(err).WithField("event_type", event.Type).Error("Failed to handle event")
				} else {
					c.logger.WithField("event_type", event.Type).Info("Event handled successfully")
				}
			} else {
				c.logger.WithField("event_type", event.Type).Warn("No handler registered for event type")
			}
		}
	}
}

func (c *Consumer) Close() error {
	c.logger.Info("Closing Kafka reader")
	return c.reader.Close()
}
