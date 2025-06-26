package producer

import (
	"context"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/events/types"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	writer *kafka.Writer
	logger *logrus.Logger
}

func NewProducer(cfg *config.KafkaConfig, logger *logrus.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topics.UserEvents,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) PublishEvent(ctx context.Context, event *types.BaseEvent) error {
	eventData, err := event.ToJSON()
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.ID.String()),
		Value: eventData,
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		p.logger.WithError(err).Error("Failed to publish event")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": event.Type,
	}).Info("Event published successfully")

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
