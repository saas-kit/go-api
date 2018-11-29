package validator

import (
	"net/http"

	"github.com/asaskevich/govalidator"
)

type (
	// Errors bag structure
	Errors struct {
		govalidator.Errors
	}
)

// ValidateStruct function
func ValidateStruct(s interface{}, r *http.Request) *Errors {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		if err2, ok := err.(govalidator.Errors); ok {
			return &Errors{err2}
		}
		e := &Errors{}
		return e
	}
	return nil
}
