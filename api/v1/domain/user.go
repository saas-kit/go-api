package domain

import (
	"fmt"
	"strings"

	"github.com/labstack/gommon/random"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// UserRepository interface
	UserRepository interface {
		GetByID(id string) (*User, error)
		GetByEmail(email string) (*User, error)
		GetAll() ([]*User, error)
		Store(*User) error
		Update(*User) error
		UpdateResetTokenTime(id string) error
		Patch(id string, data map[string]interface{}) error
		Delete(id string) error
	}

	// UserNotification interface
	UserNotification interface {
		ResetPasswordInstruction(name, email, token string) error
	}

	// User structure
	User struct {
		ID string `json:"user_id"`

		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		UserpicURL string `json:"userpic"`

		Email    string `json:"email"`
		Password string `json:"-"`
		Salt     string `json:"-"`

		Confirmed bool `json:"confirmed"`
		Disabled  bool `json:"disabled"`
	}
)

// Name helper returns full name of a customer
func (u *User) Name() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

// CheckPassword function determines the password is correct
func (u *User) CheckPassword(password string) bool {
	if u.Password == "" || u.Salt == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Salt+password))
	return err == nil
}

// SetPassword is a setter function
func (u *User) SetPassword(password string) error {
	if u.Salt == "" {
		u.Salt = random.String(16)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Salt+password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// NewUser function returns a new User object with filled data
func NewUser(email, password, firstName, lastName string) (*User, error) {
	user := &User{
		ID:        uuid.NewV1().String(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}
