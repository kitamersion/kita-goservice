package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Event struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type      string    `json:"type" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
}
