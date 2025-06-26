package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	UserCreated EventType = "user.created"
	UserUpdated EventType = "user.updated"
	UserDeleted EventType = "user.deleted"
)

type BaseEvent struct {
	ID        uuid.UUID   `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type UserCreatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
}

type UserUpdatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
}

type UserDeletedEvent struct {
	UserID uuid.UUID `json:"user_id"`
}

func NewEvent(eventType EventType, data interface{}) *BaseEvent {
	return &BaseEvent{
		ID:        uuid.New(),
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
