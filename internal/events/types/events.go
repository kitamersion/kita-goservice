package types

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Headers struct {
	ID        uuid.UUID `json:"id"`
	Timestamp string    `json:"timestamp"` // Use string for JSON serialization
}

func (e *Headers) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
