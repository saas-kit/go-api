package usecases

import (
	"errors"
	"saas-kit-api/api/v1/domain"
)

// Predefined errors
var (
	ErrWrongIvitation     = errors.New("Invitation for provided email address not found or it was expired")
	ErrInvitationNotFound = errors.New("Invitation not found")
)

type (
	// InvitationInteractor structure
	InvitationInteractor struct {
		invRepo domain.InvitationRepository
	}

	// Invitation structure
	Invitation struct {
		domain.Invitation
	}
)

// NewInvitationInteractor is factory function,
// returns a new instance of the InvitationInteractor
func NewInvitationInteractor(invRepo domain.InvitationRepository) *InvitationInteractor {
	return &InvitationInteractor{
		invRepo: invRepo,
	}
}

// GetForEmail returns invitation object
func (i *InvitationInteractor) GetForEmail(id, email string) (*Invitation, error) {
	inv, err := i.invRepo.GetByID(id)
	if err != nil || inv.IsValid(email) {
		return nil, ErrWrongIvitation
	}
	return i.invitation(inv), nil
}

// Confirm is a helper to mark invitation as confirmed
func (i *InvitationInteractor) Confirm(id string) error {
	if err := i.invRepo.Confirm(id); err != nil {
		return err
	}
	return nil
}

// Accept is a helper to mark invitation as confirmed by existed user
func (i *InvitationInteractor) Accept(id, email string) error {
	inv, err := i.invRepo.GetByID(id)
	if err != nil || inv.IsValid(email) {
		return ErrWrongIvitation
	}
	if err := i.invRepo.Confirm(id); err != nil {
		return err
	}
	return nil
}

// invitation object wrapper
func (i *InvitationInteractor) invitation(inv domain.Invitation) *Invitation {
	return &Invitation{inv}
}
