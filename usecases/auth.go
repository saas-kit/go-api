package usecases

import (
	"errors"
	"saas-kit-api/domain"
	"saas-kit-api/signeddata"
	"time"
)

const (
	sResetToken string = "reset_token"
	sSubject    string = "subject"
	sUserID     string = "user_id"
)

// Predefined errors
var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrAccessDenied       = errors.New("Access Denied")
	ErrEmailTaken         = errors.New("Provided email is already taken")
	ErrWrongIvitation     = errors.New("Invitation for provided email address not found or it was expired")
	ErrInternal           = errors.New("Internal server error. Please try again later or contact support")
	ErrWrongToken         = errors.New("Wrong reset password token")
)

type (
	// AuthInteractor structure
	AuthInteractor struct {
		userRepo      domain.UserRepository
		invRepo       domain.InvitationRepository
		projRepo      domain.ProjectRepository
		notification  domain.UserNotification
		logger        Logger
		developers    map[string]struct{}
		signingKey    string
		signedDataTTL int64
	}

	// User structure
	User struct {
		*domain.User
		IsDeveloper  bool  `json:"is_developer"`
		OriginalUser *User `json:"original_user,omitempty"`
	}
)

// SignIn use case handler
func (i *AuthInteractor) SignIn(email, password string) (*User, error) {
	user, err := i.userRepo.GetByEmail(email)
	if err != nil {
		i.logger.Error(err)
		return nil, ErrUserNotFound
	}
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}
	if err := i.userRepo.UpdateResetTokenTime(user.ID); err != nil {
		i.logger.Error(err)
		return nil, ErrInternal
	}
	return i.user(user), nil
}

// SignInAs use case handler
func (i *AuthInteractor) SignInAs(originUser *User, userID string) (*User, error) {
	if !originUser.IsDeveloper {
		return nil, ErrAccessDenied
	}
	user, err := i.userRepo.GetByID(userID)
	if err != nil {
		i.logger.Error(err)
		return nil, ErrUserNotFound
	}
	return i.user(user, originUser), nil
}

// SignOut use case handler
func (i *AuthInteractor) SignOut(userID string) error {
	if err := i.userRepo.Patch(userID, map[string]interface{}{"reset_token_at": time.Now().Unix()}); err != nil {
		i.logger.Error(err)
	}
	return nil
}

// SignUp use case handler
func (i *AuthInteractor) SignUp(email, password, firstName, lastName string) (*User, error) {
	if _, err := i.userRepo.GetByEmail(email); err == nil {
		i.logger.Error(err)
		return nil, ErrEmailTaken
	}
	user, err := domain.NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}
	if err := i.userRepo.Store(user); err != nil {
		return nil, err
	}
	return i.user(user), nil
}

// SignUpViaInvitation use case handler
func (i *AuthInteractor) SignUpViaInvitation(email, password, firstName, lastName, invitationID string) (*User, error) {
	inv, err := i.invRepo.GetByID(invitationID)
	if err != nil || inv.CheckEmail(email) || inv.Project.Disabled {
		i.logger.Error(err)
		return nil, ErrWrongIvitation
	}
	if _, err = i.userRepo.GetByEmail(email); err == nil {
		i.logger.Error(err)
		return nil, ErrEmailTaken
	}
	user, err := domain.NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}
	if err = i.userRepo.Store(user); err != nil {
		return nil, err
	}
	if err = i.invRepo.Confirm(invitationID); err != nil {
		if err2 := i.userRepo.Delete(user.ID); err2 != nil {
			i.logger.Error(err2)
		}
		i.logger.Error(err)
		return nil, ErrInternal
	}
	if i.projRepo.AddMember(inv.Project.ID, user.ID, inv.Role); err != nil {
		i.logger.Error(err)
		return nil, ErrInternal
	}
	return i.user(user), nil
}

// ForgotPassword use case handler
func (i *AuthInteractor) ForgotPassword(email string) error {
	user, err := i.userRepo.GetByEmail(email)
	if err == nil {
		i.logger.Error(err)
		return ErrUserNotFound
	}
	signedData, err := signeddata.Encode(i.signingKey, map[string]interface{}{
		sUserID:  user.ID,
		sSubject: sResetToken,
	}, time.Now().Add(time.Duration(i.signedDataTTL)*time.Second))
	if err != nil {
		return err
	}
	if err = i.notification.ResetPasswordInstruction(user.Name(), user.Email, signedData); err != nil {
		i.logger.Error(err)
		return ErrInternal
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
	user, err := i.userRepo.GetByID(userID.(string))
	if err != nil {
		return err
	}
	user.SetPassword(password)
	if err := i.userRepo.Update(user); err != nil {
		return err
	}
	return nil
}

// user function is a helper to the composition a new User object based on the domain.User object
func (i *AuthInteractor) user(user *domain.User, originUser ...*User) *User {
	u := &User{user, false, nil}
	u.IsDeveloper = i.isDeveloper(user.Email)
	if len(originUser) > 0 {
		u.OriginalUser = originUser[0]
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
