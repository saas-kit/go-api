package usecases

import "saas-kit-api/app/auth/domain"

type (
	// ProfileInteractor structure
	ProfileInteractor struct {
		userRepo domain.UserRepository
		logger   Logger
	}
)

// NewProfileInteractor is a factory function,
// returns a new instance of the ProfileInteractor
func NewProfileInteractor(userRepo domain.UserRepository, logger Logger) *ProfileInteractor {
	return &ProfileInteractor{
		userRepo: userRepo,
		logger:   logger,
	}
}

// ChangePassword use case handler
func (i *ProfileInteractor) ChangePassword(old, new string) error {
	return nil
}

// ChangeEmail use case handler
func (i *ProfileInteractor) ChangeEmail(email string) error {
	return nil
}

// ConfirmEmail use case handler
func (i *ProfileInteractor) ConfirmEmail(token string) error {
	return nil
}

// ResendEmailConfirmation use case handler
func (i *ProfileInteractor) ResendEmailConfirmation(userID string) error {
	return nil
}
