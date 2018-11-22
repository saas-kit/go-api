package usecases

import (
	"errors"
	"saas-kit-api/api/v1/domain"
)

// Predefined errors
var (
	ErrUserNotFound       = errors.New("User not found")
	ErrEmailTaken         = errors.New("Provided email is already taken")
	ErrCouldNotDeleteUser = errors.New("Could not delete user")
)

type (
	// User structure
	User struct {
		domain.User
		IsDeveloper  bool  `json:"is_developer"`
		OriginalUser *User `json:"original_user,omitempty"`
	}

	// UserInteractor structure
	UserInteractor struct {
		userRepo   domain.UserRepository
		developers map[string]struct{}
	}
)

// NewUserInteractor is a factory function,
// returns a new instance of the UserInteractor structure
func NewUserInteractor(ur domain.UserRepository, developers map[string]struct{}) *UserInteractor {
	return &UserInteractor{
		userRepo:   ur,
		developers: developers,
	}
}

// SetOrigin helper to set original user
func (u *User) SetOrigin(user *User) {
	u.OriginalUser = user
}

// GetByID to fetch user by id
func (i *UserInteractor) GetByID(id string) (*User, error) {
	user, err := i.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return i.user(user), nil
}

// GetByEmail usecase handler
func (i *UserInteractor) GetByEmail(email string) (*User, error) {
	user, err := i.userRepo.GetByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return i.user(user), nil
}

// Create new user
func (i *UserInteractor) Create(email, password, firstName, lastName string) (*User, error) {
	if _, err := i.userRepo.GetByEmail(email); err == nil {
		return nil, ErrEmailTaken
	}
	user, err := domain.NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}
	if err := i.userRepo.Store(&user); err != nil {
		return nil, err
	}
	return i.user(user), nil
}

// UpdatePassword usecase handler
func (i *UserInteractor) UpdatePassword(id, password string) error {
	user, err := i.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	user.SetPassword(password)
	if err := i.userRepo.Update(&user); err != nil {
		return err
	}
	return nil
}

// UpdateResetTokenTime usecase handler
func (i *UserInteractor) UpdateResetTokenTime(id string) error {
	return i.userRepo.UpdateResetTokenTime(id)
}

// Delete user by id
func (i *UserInteractor) Delete(id string) error {
	if i.userRepo.Delete(id) != nil {
		return ErrCouldNotDeleteUser
	}
	return nil
}

// user function is a helper to the composition a new User object based on the domain.User object
func (i *UserInteractor) user(user domain.User, originUser ...*User) *User {
	u := &User{user, false, nil}
	u.IsDeveloper = i.isDeveloper(user.Email)
	if len(originUser) > 0 {
		u.OriginalUser = originUser[0]
	}
	return u
}

// isDeveloper helper determines email address belongs to developer user
func (i *UserInteractor) isDeveloper(email string) bool {
	if _, ok := i.developers[email]; ok {
		return true
	}
	return false
}
