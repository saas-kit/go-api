package usecases

import (
	"errors"
	"saas-kit-api/pkg/signeddata"
	"time"

	"saas-kit-api/api/v1/domain"
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
	// AuthInteractor structure
	AuthInteractor struct {
		user          authUserInteractor
		invitation    authInvitationInteractor
		project       authProjectInteractor
		notification  domain.UserNotification
		signingKey    string
		signedDataTTL int64
	}

	authUserInteractor interface {
		GetByID(id string) (*User, error)
		GetByEmail(email string) (*User, error)
		Create(email, password, firstName, lastName string) (*User, error)
		UpdatePassword(id, password string) error
		UpdateResetTokenTime(id string) error
		Delete(id string) error
	}

	authInvitationInteractor interface {
		GetForEmail(id, email string) (*Invitation, error)
		Confirm(id string) error
	}

	authProjectInteractor interface {
		AddMember(projectID, userID string, role domain.Role) error
	}
)

// NewAuthInteractor is a factory function,
// returns a new instance of the AuthInteractor
func NewAuthInteractor(
	user authUserInteractor,
	invitation authInvitationInteractor,
	project authProjectInteractor,
	notification domain.UserNotification,
	signingKey string,
	signedDataTTL int64,
) *AuthInteractor {
	return &AuthInteractor{
		user:          user,
		invitation:    invitation,
		project:       project,
		notification:  notification,
		signingKey:    signingKey,
		signedDataTTL: signedDataTTL,
	}
}

// SignIn use case handler
func (i *AuthInteractor) SignIn(email, password string) (*User, error) {
	user, err := i.user.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}
	if err := i.user.UpdateResetTokenTime(user.ID); err != nil {
		return nil, err
	}
	return user, nil
}

// SignInAs use case handler
func (i *AuthInteractor) SignInAs(originUser *User, userID string) (*User, error) {
	if !originUser.IsDeveloper {
		return nil, ErrAccessDenied
	}
	user, err := i.user.GetByID(userID)
	if err != nil {
		return nil, err
	}
	user.SetOrigin(originUser)
	return user, nil
}

// SignOut use case handler
func (i *AuthInteractor) SignOut(userID string) error {
	return i.user.UpdateResetTokenTime(userID)
}

// SignUp use case handler
func (i *AuthInteractor) SignUp(email, password, firstName, lastName string) (*User, error) {
	user, err := i.user.Create(email, password, firstName, lastName)
	if err == nil {
		return nil, err
	}
	return user, nil
}

// SignUpViaInvitation use case handler
func (i *AuthInteractor) SignUpViaInvitation(email, password, firstName, lastName, invitationID string) (*User, error) {
	inv, err := i.invitation.GetForEmail(invitationID, email)
	if err != nil {
		return nil, err
	}
	user, err := i.user.Create(email, password, firstName, lastName)
	if err == nil {
		return nil, err
	}
	if err = i.invitation.Confirm(invitationID); err != nil {
		i.user.Delete(user.ID)
		return nil, err
	}
	if err := i.project.AddMember(inv.Project.ID, user.ID, inv.Role); err != nil {
		return nil, err
	}
	return user, nil
}

// ForgotPassword use case handler
func (i *AuthInteractor) ForgotPassword(email string) error {
	user, err := i.user.GetByEmail(email)
	if err == nil {
		return err
	}
	signedData, err := signeddata.Encode(i.signingKey, map[string]interface{}{
		sUserID:  user.ID,
		sSubject: sResetToken,
	}, time.Now().Add(time.Duration(i.signedDataTTL)*time.Second))
	if err != nil {
		return err
	}
	if err = i.notification.ResetPasswordInstruction(user.Name(), user.Email, signedData); err != nil {
		return err
	}
	return nil
}

// ResetPassword use case handler
func (i *AuthInteractor) ResetPassword(token, password string) error {
	payload, err := signeddata.Decode(i.signingKey, token)
	if err != nil {
		return err
	}
	if val, ok := payload[sSubject]; !ok || val != sResetToken {
		return ErrWrongToken
	}
	userID, ok := payload[sUserID]
	if !ok {
		return ErrWrongToken
	}
	if err := i.user.UpdatePassword(userID.(string), password); err != nil {
		return err
	}
	return nil
}
