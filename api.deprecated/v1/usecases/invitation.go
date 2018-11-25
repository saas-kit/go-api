package usecases

import (
	"errors"
	"saas-kit-api/api/v1/domain"
)

// Predefined errors
var (
	ErrWrongIvitation     = errors.New("Invitation for provided email address not found or it was expired")
	ErrInvitationNotFound = errors.New("Invitation not found or it was accepted")
	ErrProjectNotFound    = errors.New("Project not found or it is disabled")
)

type (
	// InvitationInteractor structure
	InvitationInteractor struct {
		invRepo      domain.InvitationRepository
		notification domain.UserNotification
		project      invProjectInteractor
	}

	// Invitation structure
	Invitation struct {
		*domain.Invitation
		Project *Project
	}

	invProjectInteractor interface {
		GetByID(id string) (*Project, error)
	}
)

// NewInvitationInteractor is factory function,
// returns a new instance of the InvitationInteractor
func NewInvitationInteractor(invRepo domain.InvitationRepository, notification domain.UserNotification) *InvitationInteractor {
	return &InvitationInteractor{
		invRepo:      invRepo,
		notification: notification,
	}
}

// GetForEmail returns invitation object
func (i *InvitationInteractor) GetForEmail(id, email string) (*Invitation, error) {
	inv, err := i.invRepo.GetByID(id)
	if err != nil || inv.CheckEmail(email) {
		return nil, ErrWrongIvitation
	}
	return i.wrap(inv), nil
}

// GetListForProject returns invitations list by project id
func (i *InvitationInteractor) GetListForProject(projectID string) ([]Invitation, error) {
	invs, err := i.invRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return i.wrapList(invs), nil
}

// GetListForUser returns invitations list by user email
func (i *InvitationInteractor) GetListForUser(email string) ([]Invitation, error) {
	invs, err := i.invRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return i.wrapList(invs), nil
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
	if err != nil || inv.CheckEmail(email) {
		return ErrWrongIvitation
	}
	if err := i.invRepo.Confirm(id); err != nil {
		return err
	}
	return nil
}

// Sent a new invitation
func (i *InvitationInteractor) Sent(projectID, email string, role domain.Role) (*Invitation, error) {
	project, err := i.project.GetByID(projectID)
	if err != nil || !project.Disabled {
		return nil, ErrProjectNotFound
	}
	inv := domain.NewInvitation(projectID, email, role)
	if err := i.invRepo.Store(inv); err != nil {
		return nil, err
	}
	return i.wrap(inv), nil
}

// Remove invitation
func (i *InvitationInteractor) Remove(id string) error {
	inv, err := i.invRepo.GetByID(id)
	if err != nil || inv.Confirmed {
		return ErrInvitationNotFound
	}
	return i.invRepo.Delete(id)
}

// invitation object wrapper
func (i *InvitationInteractor) wrap(inv *domain.Invitation) *Invitation {
	return &Invitation{inv, nil}
}

// wrapper for domain.Project structures list
func (i *InvitationInteractor) wrapList(invs []domain.Invitation) []Invitation {
	list := make([]Invitation, 0, len(invs))
	for _, i := range invs {
		list = append(list, Invitation{&i, nil})
	}
	return list
}
