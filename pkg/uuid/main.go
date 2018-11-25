package uuid

import (
	"github.com/satori/go.uuid"
)

type (
	// UUID structure
	UUID struct {
		uuid uuid.UUID
	}
)

// String returns generated uuid v1 string
func (u UUID) String() string {
	return u.uuid.String()
}

// NewV1 creates a new UUID v1 object
func NewV1() UUID {
	return UUID{uuid.NewV1()}
}
