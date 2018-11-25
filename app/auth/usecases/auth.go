package usecases

import (
	"errors"
	"saas-kit-api/app/auth/domain"
)

const (
	sResetToken string = "reset_token"
	sSubject    string = "subject"
	sUserID     string = "user_id"
)

// Predefined errors
var (
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrAccessDenied       = errors.New("Access Denied")
	ErrWrongToken         = errors.New("Wrong reset password token")
)

type (
	// AuthInteractor struct
	AuthInteractor struct {
		userRepo      domain.UserRepository
		signingKey    string
		signedDataTTL int64
		developers    map[string]struct{}
		logger        Logger
	}

	// User structure
	User struct {
		*domain.User
		IsDeveloper bool  `json:"is_developer"`
		OriginUser  *User `json:"origin_user"`
	}
)

// NewAuthInteractor is a factory function,
// returns a new instance of the AuthInteractor
func NewAuthInteractor(
	userRepo domain.UserRepository,
	signingKey string, signedDataTTL int64,
	developers map[string]struct{},
	logger Logger,
) *AuthInteractor {
	return &AuthInteractor{
		userRepo:      userRepo,
		signingKey:    signingKey,
		signedDataTTL: signedDataTTL,
		developers:    developers,
		logger:        logger,
	}
}

// SignIn use case handler
func (i *AuthInteractor) SignIn(email, password string) (*User, error) {
	return nil, nil
}

// SignInAs use case handler
func (i *AuthInteractor) SignInAs(userID, originID string) (*User, error) {
	return nil, nil
}

// SignUp use case handler
func (i *AuthInteractor) SignUp(email, password string) (*User, error) {
	return nil, nil
}

// SignOut use case handler
func (i *AuthInteractor) SignOut(userID string) error {
	return nil
}

// BackToOrigin use case handler
func (i *AuthInteractor) BackToOrigin(userID string) (*User, error) {
	return nil, nil
}

// ForgotPassword use case handler
func (i *AuthInteractor) ForgotPassword(email string) error {
	return nil
}

// ResetPassword use case handler
func (i *AuthInteractor) ResetPassword(token, newPassword string) error {
	return nil
}

// user function is a helper to the composition a new User object based on the domain.User object
func (i *AuthInteractor) wrap(user *domain.User, originUser ...*User) *User {
	u := &User{user, false, nil}
	u.IsDeveloper = i.isDeveloper(user.Email)
	if len(originUser) > 0 {
		u.OriginUser = originUser[0]
	}
	return u
}

// isDeveloper helper determines email address belongs to developer user
func (i *AuthInteractor) isDeveloper(email string) bool {
	if _, ok := i.developers[email]; ok {
		return true
	}
	return false
}
