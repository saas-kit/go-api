package http

import (
	"github.com/labstack/echo"

	"saas-kit-api/api/v1/usecases"
)

type (
	// AuthWebService structure
	AuthWebService struct {
		authInteractor authInteractor
	}

	// AuthInteractor interface
	authInteractor interface {
		SignIn(email, password string) (usecases.User, error)
		SignInAs(originUser usecases.User, userID string) (usecases.User, error)
		SignOut(userID string) error
		SignUp(email, password, firstName, lastName string) (usecases.User, error)
		SignUpViaInvitation(email, password, firstName, lastName, invitationID string) (usecases.User, error)
		ForgotPassword(email string) error
		ResetPassword(token, password string) error
	}
)

// NewAuthWebService is a factory function,
// return a new instance of the AuthWebService
func NewAuthWebService(ai authInteractor) *AuthWebService {
	return &AuthWebService{
		authInteractor: ai,
	}
}

// SetUp auth web service
func (ws *AuthWebService) SetUp(router *echo.Group) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signin", ws.SignIn)
		authGroup.POST("/signinas", ws.SignInAs)
		authGroup.POST("/signout", ws.SignOut)
		authGroup.POST("/signup", ws.SignUp)
		authGroup.POST("/forgot-password", ws.ForgotPassword)
		authGroup.POST("/reset-password", ws.ResetPassword)
	}
}

// SignIn route handler
func (ws *AuthWebService) SignIn(c echo.Context) error {
	return nil
}

// SignInAs route handler
func (ws *AuthWebService) SignInAs(c echo.Context) error {
	return nil
}

// SignOut route handler
func (ws *AuthWebService) SignOut(c echo.Context) error {
	return nil
}

// SignUp route handler
func (ws *AuthWebService) SignUp(c echo.Context) error {
	return nil
}

// SignUpViaInvitation route handler
func (ws *AuthWebService) SignUpViaInvitation(c echo.Context) error {
	return nil
}

// ForgotPassword route handler
func (ws *AuthWebService) ForgotPassword(c echo.Context) error {
	return nil
}

// ResetPassword route handler
func (ws *AuthWebService) ResetPassword(c echo.Context) error {
	return nil
}
