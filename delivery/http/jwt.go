package http

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	// JWTContextKey is a key to save user data into context
	JWTContextKey string = "JWT"
)

// JWTClaims is a custom JWT claims structure
type JWTClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// JWT Middleware
func JWT(signingKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		BeforeFunc:     beforeFunc,
		SuccessHandler: successHandler,
		ErrorHandler:   errorHandler,
		SigningKey:     []byte(signingKey),
		SigningMethod:  middleware.AlgorithmHS256,
		ContextKey:     JWTContextKey,
		Claims:         jwt.MapClaims{},
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
