package domain

import (
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type (
	// Invitation structure
	Invitation struct {
		ID        string    `json:"invitation_id"`
		ProjectID string    `json:"project_id"`
		Email     string    `json:"email"`
		Project   Project   `json:"project"`
		Role      Role      `json:"role"`
		SentAt    time.Time `json:"sent_at"`
		Confirmed bool      `json:"confirmed"`
	}

	// InvitationRepository interface
	InvitationRepository interface {
		GetByID(id string) (*Invitation, error)
		GetByEmail(email string) ([]Invitation, error)
		GetByProject(projectID string) ([]Invitation, error)
		GetList(limit, offset int) ([]Invitation, error)
		Store(*Invitation) error
		Update(*Invitation) error
		Delete(id string) error
		Confirm(id string) error
	}
)

// NewInvitation is a factory function to create a new instance of the Invitation structure
func NewInvitation(projectID, email string, role Role) *Invitation {
	return &Invitation{
		ID:        uuid.NewV1().String(),
		ProjectID: projectID,
		Email:     email,
		Role:      role,
		SentAt:    time.Now(),
	}
}

// CheckEmail checks if the provided email address equals email address from an invitation
func (i *Invitation) CheckEmail(email string) bool {
	return strings.ToLower(email) == strings.ToLower(i.Email)
}

// IsValid checks if an invitation is valid
func (i *Invitation) IsValid(email string) bool {
	return i.CheckEmail(email) && !i.Project.Disabled
}
