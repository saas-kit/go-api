package domain

import (
	"fmt"
	"saas-kit-api/pkg/uuid"
)

type (
	// Address struct
	Address struct {
		ID      string `json:"id" db:"id"`
		Line1   string `json:"line_1" db:"line_1"`
		Line2   string `json:"line_2" db:"line_2"`
		City    string `json:"city" db:"city"`
		State   string `json:"state" db:"state"`
		Country string `json:"country" db:"country"`
		ZipCode string `json:"zip_code" db:"zip_code"`
	}

	// AddressRepository interface
	AddressRepository interface {
		GetByID(id string) (*Address, error)
		Store(*Address) error
		Update(*Address) error
		Delete(id string) error
		ForceDelete(id string) error
	}
)

func (a *Address) String() string {
	return fmt.Sprintf("%s %s, %s, %s %s, %s", a.Line1, a.Line2, a.City, a.State, a.ZipCode, a.Country)
}

// NewAddress is a new Address factory
func NewAddress(line1, line2, city, state, country, zip string) *Address {
	return &Address{
		ID:      uuid.NewV1().String(),
		Line1:   line1,
		Line2:   line2,
		City:    city,
		State:   state,
		Country: country,
		ZipCode: zip,
	}
}
