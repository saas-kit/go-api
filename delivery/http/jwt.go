package http

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/satori/go.uuid"
)

const (
	// JWTContextKey is a key to save user data into context
	JWTContextKey string = "JWT"
)

type (
	// JWT structure
	JWT struct {
		SigningKey []byte
		TTL        int64
	}

	// JWTResponse structure
	JWTResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}

	// JWTClaims is a custom JWT claims structure
	JWTClaims struct {
		jwt.StandardClaims
		ID             string `json:"id"`
		Name           string `json:"name"`
		Role           string `json:"role"`
		OriginalUserID string `json:"original_user_id,omitempty"`
		Confirmed      bool   `json:"confirmed"`
	}

	// User claims interface
	userJWTClaims interface {
		GetID() string
		GetName() string
		GetRole() string
		IsConfirmed() bool
	}

	// DefaultUserJWTClaims structure
	DefaultUserJWTClaims struct {
		ID        string
		Name      string
		Role      string
		Confirmed bool
	}
)

// Middleware function to handle JWT from request
func (j *JWT) Middleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		BeforeFunc:     beforeFunc,
		SuccessHandler: successHandler,
		ErrorHandler:   errorHandler,
		SigningKey:     j.SigningKey,
		SigningMethod:  middleware.AlgorithmHS256,
		ContextKey:     JWTContextKey,
		Claims:         JWTClaims{},
		TokenLookup:    "header:" + echo.HeaderAuthorization,
		AuthScheme:     "Bearer",
	})
}

// Function runs before JWT middleware
func beforeFunc(c echo.Context) {}

// Successfull handler
func successHandler(c echo.Context) {}

// Error handler
func errorHandler(err error) error {
	if err == middleware.ErrJWTMissing {
		return NewHTTPError(http.StatusUnauthorized, middleware.ErrJWTMissing.Message)
	}
	if err, ok := err.(*echo.HTTPError); ok {
		return NewHTTPError(err.Code, err.Message)
	}
	return NewHTTPError(http.StatusUnauthorized, err.Error())
}

// NewToken returns a new JWT object
func (j *JWT) NewToken(claims JWTClaims) (JWTResponse, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.SigningKey)
	if err != nil {
		return JWTResponse{}, err
	}
	return JWTResponse{
		AccessToken: token,
		ExpiresAt:   claims.ExpiresAt,
	}, nil
}

// NewClaims returns new JWTClaims structure instance with populated default fields
func (j *JWT) NewClaims(user userJWTClaims, originalID ...string) JWTClaims {
	claims := JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(j.TTL) * time.Second).Unix(),
			Id:        uuid.NewV1().String(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
		user.GetID(),
		user.GetName(),
		user.GetRole(),
		"",
		user.IsConfirmed(),
	}
	// workaround to make optional parameter originalID
	if len(originalID) > 0 {
		claims.OriginalUserID = originalID[0]
	}
	return claims
}

// GetID returns user id
func (u DefaultUserJWTClaims) GetID() string {
	return u.ID
}

// GetName returns user name
func (u DefaultUserJWTClaims) GetName() string {
	return u.Name
}

// GetRole returns user role
func (u DefaultUserJWTClaims) GetRole() string {
	return u.Role
}

// IsConfirmed returns user confirmation status
func (u DefaultUserJWTClaims) IsConfirmed() bool {
	return u.Confirmed
}
