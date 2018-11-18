package http

import (
	"github.com/labstack/echo"
)

// Default API routes
func setupRoutes(r *echo.Group, cnf routesConfig) {
	r.Use(JSONHeadersMiddleware)
}
