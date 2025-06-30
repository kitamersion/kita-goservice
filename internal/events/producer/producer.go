package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/events/types"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	writer *kafka.Writer
	logger *logrus.Logger
}

func NewProducer(cfg *config.KafkaConfig, logger *logrus.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topics.UserEvents,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // Strong durability
		BatchSize:    1,                // default is fine, or increase for higher throughput
		BatchTimeout: 10 * time.Millisecond,
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) PublishEvent(ctx context.Context, protoEvent proto.Message) error {
	// Serialize the proto event to JSON
	serializedEvent, err := protojson.Marshal(protoEvent)
	if err != nil {
		p.logger.WithError(err).Error("Failed to serialize proto event")
		return fmt.Errorf("failed to serialize proto event: %w", err)
	}

	headers := types.Headers{
		ID:        uuid.New(),
		Timestamp: fmt.Sprint(time.Now().Unix()), // Use Unix timestamp as string for JSON serialization
	}

	// Get the event type from the proto message
	eventType := string(protoEvent.ProtoReflect().Descriptor().FullName())

	message := kafka.Message{
		Key: []byte(headers.ID.String()),
		Headers: []kafka.Header{
			{Key: "event_id", Value: []byte(headers.ID.String())},
			{Key: "timestamp", Value: []byte(headers.Timestamp)},
			{Key: "event_type", Value: []byte(eventType)},
			{Key: "content_type", Value: []byte("application/json")},
		},
		Value: serializedEvent,
	}

	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		p.logger.WithError(err).Error("Failed to publish event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.WithFields(logrus.Fields{
		"event_id":   headers.ID.String(),
		"event_type": eventType,
	}).Info("Event published successfully")

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
