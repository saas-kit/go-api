package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

// Predefined errors
var (
	ErrMalformedJSON = echo.NewHTTPError(http.StatusBadRequest, "Malformed JSON")
)

type (
	// Binder structure
	Binder struct{}

	// ValidatableStruct interface
	ValidatableStruct interface {
		ValidationRules() map[string][]string
		ValidationMessages() map[string][]string
	}
)

// Bind is a custom implementation of the `Binder#Bind` function.
func (cb *Binder) Bind(i interface{}, c echo.Context) error {
	if data, ok := i.(ValidatableStruct); ok {
		v := govalidator.New(govalidator.Options{
			Request:  c.Request(),
			Data:     data,
			Rules:    data.ValidationRules(),
			Messages: data.ValidationMessages(),
		})
		if err := v.ValidateJSON(); len(err) > 0 {
			if _, ok := err["_error"]; ok {
				return ErrMalformedJSON
			}
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	} else {
		db := new(echo.DefaultBinder)
		if err := db.Bind(i, c); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	return nil
}
