package domain

import (
	"saas-kit-api/pkg/hash"
	"saas-kit-api/pkg/random"
	"saas-kit-api/pkg/uuid"
	"time"
)

type (
	// User structure
	User struct {
		ID string `json:"id"`

		Email        string     `json:"email"`
		Password     string     `json:"-"`
		Salt         string     `json:"-"`
		ResetTokenAt *time.Time `json:"-"`

		Confirmed bool `json:"confirmed"`
		Disabled  bool `json:"disabled"`
	}

	// UserRepository interface
	UserRepository interface {
		GetByID(id string) (*User, error)
		GetByEmail(email string) (*User, error)
		Store(*User) error
		ChangeEmail(id, email string) error
		UpdatePassword(id, password string) error
		UpdateResetTokenTime(id string) error
		Confirm(id string) error
		Unconfirm(id string) error
		Disable(id string) error
		Enable(id string) error
		Delete(id string) error
	}
)

// CheckPassword function determines the password is correct
func (u *User) CheckPassword(password string) bool {
	if u.Password == "" || u.Salt == "" {
		return false
	}
	err := hash.Compare(u.Password, u.Salt+password)
	return err == nil
}

// SetPassword is a setter function
func (u *User) SetPassword(password string) error {
	if u.Salt == "" {
		u.Salt = random.String(16)
	}
	hashStr, err := hash.New(u.Salt + password)
	if err != nil {
		return err
	}
	u.Password = hashStr
	return nil
}

// NewUser function returns a new User object with filled data
func NewUser(email, password string) (*User, error) {
	user := &User{
		ID:    uuid.NewV1().String(),
		Email: email,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}
