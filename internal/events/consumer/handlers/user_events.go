package handlers

import (
	"context"
	"fmt"

	"github.com/kitamersion/go-goservice/internal/events/proto/events/userpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserEventHandlers struct {
	logger *logrus.Logger
}

func NewUserEventHandlers(logger *logrus.Logger) *UserEventHandlers {
	return &UserEventHandlers{
		logger: logger,
	}
}

func (h *UserEventHandlers) HandleUserCreated(ctx context.Context, eventType string, headers map[string]string, payload []byte) error {
	var event userpb.UserCreated
	if err := protojson.Unmarshal(payload, &event); err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal UserCreated event")
		return fmt.Errorf("failed to unmarshal UserCreated event: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"event_id":   headers["event_id"],
		"event_type": eventType,
		"user_id":    event.Id,
		"email":      event.Email,
		"name":       event.Name,
	}).Info("User created event processed")

	// Add your business logic here
	// For example: send welcome email, update cache, etc.

	return nil
}

func (h *UserEventHandlers) HandleUserUpdated(ctx context.Context, eventType string, headers map[string]string, payload []byte) error {
	var event userpb.UserUpdated
	if err := protojson.Unmarshal(payload, &event); err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal UserUpdated event")
		return fmt.Errorf("failed to unmarshal UserUpdated event: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"event_id":   headers["event_id"],
		"event_type": eventType,
		"user_id":    event.Id,
	}).Info("User updated event processed")

	// Add your business logic here
	// For example: update cache, sync with external systems, etc.

	return nil
}

func (h *UserEventHandlers) HandleUserDeleted(ctx context.Context, eventType string, headers map[string]string, payload []byte) error {
	var event userpb.UserDeleted
	if err := protojson.Unmarshal(payload, &event); err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal UserDeleted event")
		return fmt.Errorf("failed to unmarshal UserDeleted event: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"event_id":   headers["event_id"],
		"event_type": eventType,
		"user_id":    event.Id,
	}).Info("User deleted event processed")

	// Add your business logic here
	// For example: cleanup user data, update search index, etc.

	return nil
}
