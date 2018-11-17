package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Default API routes
func setupRoutes(r *echo.Group, cnf routesConfig) {
	jwt := JWT{cnf.JWTSigningKey(), cnf.JWTTTL()}
	r.Use(JSONHeadersMiddleware)
	r.GET("/test", testHandler(jwt))
	r.GET("/restricted", testHandler(jwt), jwt.Middleware())
	r.POST("/validation", postHandler)
}

func testHandler(jwt JWT) echo.HandlerFunc {
	user := DefaultUserJWTClaims{
		ID:        "userid",
		Name:      "John Doe",
		Role:      "customer",
		Confirmed: true,
	}
	return func(c echo.Context) error {
		token, err := jwt.NewToken(jwt.NewClaims(user))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":      http.StatusText(http.StatusOK),
			"server_time": time.Now(),
			"jwt":         token,
		})
	}
}

// UserRequest struct
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ValidationRules func
func (r UserRequest) ValidationRules() map[string][]string {
	return map[string][]string{
		"email": []string{"required", "email"},
		"name":  []string{"required", "min:5", "max:50"},
	}
}

// ValidationMessages func
func (r UserRequest) ValidationMessages() map[string][]string {
	return map[string][]string{
		"email": []string{
			"required:email address is required field",
			"email:it must be valid email address",
		},
		"name": []string{
			"required:username is required",
			"min:username must be at least 5 chars",
			"max:username must be less then 50 chars",
		},
	}
}

func postHandler(c echo.Context) error {
	form := &UserRequest{}
	if err := c.Bind(form); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, form)
}
