package http

import "github.com/labstack/echo"

// Health check endpoint
func healthCheckHandler(c echo.Context) error {
	return c.String(200, "OK")
}
