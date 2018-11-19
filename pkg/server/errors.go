package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type (
	// Error is a default HTTP error structure
	Error struct {
		Code   int         `json:"code"`
		Title  string      `json:"title"`
		Detail interface{} `json:"detail,omitempty"`
	}
)

// Error makes it compatible with `error` interface.
func (e *Error) Error() string {
	return fmt.Sprintf("code=%d, title=%s, detail=%v", e.Code, e.Title, e.Detail)
}

// NewHTTPError halper to return formated HTTP error object
func NewHTTPError(code int, message ...interface{}) *Error {
	e := &Error{
		Code:  code,
		Title: http.StatusText(code),
	}
	if counter := len(message); counter > 0 {
		if counter == 1 {
			e.Detail = message[0]
		} else {
			e.Detail = message
		}
	}
	return e
}

// Default HTTP errors handler
func defaultHTTPErrorHandler(err error, c echo.Context) {
	var e *Error
	if he, ok := err.(*echo.HTTPError); ok {
		e = NewHTTPError(he.Code, he.Message)
		if he.Internal != nil {
			c.Logger().Error(he.Internal)
		}
	} else if de, ok := err.(*Error); ok {
		e = de
	} else {
		c.Logger().Error(err)
		e = NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := c.JSON(e.Code, map[string]interface{}{"error": e}); err != nil {
		panic(err)
	}
}
