package handlers

import (
	"context"
	"encoding/json"

	"github.com/kitamersion/go-goservice/internal/events/types"
	"github.com/sirupsen/logrus"
)

type UserEventHandlers struct {
	logger *logrus.Logger
}

func NewUserEventHandlers(logger *logrus.Logger) *UserEventHandlers {
	return &UserEventHandlers{
		logger: logger,
	}
}

func (h *UserEventHandlers) HandleUserCreated(ctx context.Context, event *types.BaseEvent) error {
	var userCreatedEvent types.UserCreatedEvent

	// Convert the Data interface{} back to the specific event type
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(dataBytes, &userCreatedEvent); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"user_id":  userCreatedEvent.UserID,
		"email":    userCreatedEvent.Email,
		"name":     userCreatedEvent.Name,
	}).Info("User created event processed")

	// Add your business logic here
	// For example: send welcome email, update cache, etc.

	return nil
}

func (h *UserEventHandlers) HandleUserUpdated(ctx context.Context, event *types.BaseEvent) error {
	var userUpdatedEvent types.UserUpdatedEvent

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(dataBytes, &userUpdatedEvent); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"user_id":  userUpdatedEvent.UserID,
		"email":    userUpdatedEvent.Email,
		"name":     userUpdatedEvent.Name,
	}).Info("User updated event processed")

	// Add your business logic here

	return nil
}

func (h *UserEventHandlers) HandleUserDeleted(ctx context.Context, event *types.BaseEvent) error {
	var userDeletedEvent types.UserDeletedEvent

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(dataBytes, &userDeletedEvent); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"user_id":  userDeletedEvent.UserID,
	}).Info("User deleted event processed")

	// Add your business logic here
	// For example: cleanup user data, update search index, etc.

	return nil
}
