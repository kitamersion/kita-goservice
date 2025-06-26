package consumer

import (
	"context"
	"encoding/json"

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
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topics.UserEvents,
		GroupID: "user-service-consumer",
	})

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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.WithError(err).Error("Failed to read message")
				continue
			}

			var event types.BaseEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				c.logger.WithError(err).Error("Failed to unmarshal event")
				continue
			}

			if handler, exists := c.handlers[event.Type]; exists {
				if err := handler(ctx, &event); err != nil {
					c.logger.WithError(err).WithField("event_type", event.Type).Error("Failed to handle event")
				}
			} else {
				c.logger.WithField("event_type", event.Type).Warn("No handler registered for event type")
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
