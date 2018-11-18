package domain

import (
	"strings"
	"time"
)

type (
	// Invitation structure
	Invitation struct {
		ID        string    `json:"invitation_id"`
		Email     string    `json:"email"`
		Project   Project   `json:"project"`
		Role      Role      `json:"role"`
		SentAt    time.Time `json:"sent_at"`
		Confirmed bool      `json:"confirmed"`
	}

	// InvitationRepository interface
	InvitationRepository interface {
		GetByID(id string) (Invitation, error)
		GetByEmail(email string) (Invitation, error)
		GetByProject(projectID string) (Invitation, error)
		GetAll() ([]Invitation, error)
		Store(*Invitation) error
		Update(*Invitation) error
		Patch(id string, data map[string]interface{}) error
		Delete(id string) error
		Confirm(id string) error
	}
)

// CheckEmail checks if the provided email address equals email address from an invitation
func (i *Invitation) CheckEmail(email string) bool {
	return strings.ToLower(email) == strings.ToLower(i.Email)
}
