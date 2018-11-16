package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// JSONHeadersMiddleware checks content type and accept headers
func JSONHeadersMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != echo.GET {
			cth := strings.ToLower(c.Request().Header.Get(echo.HeaderContentType))
			if !strings.Contains(cth, echo.MIMEApplicationJSON) {
				return NewHTTPError(http.StatusUnsupportedMediaType)
			}
		}
		if !strings.Contains(c.Request().Header.Get(echo.HeaderAccept), echo.MIMEApplicationJSON) {
			return NewHTTPError(http.StatusNotAcceptable)
		}
		return next(c)
	}
}
