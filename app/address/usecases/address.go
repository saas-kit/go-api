package usecases

import (
	"saas-kit-api/app/address/domain"
)

type (
	// AddressInteractor struct
	AddressInteractor struct {
		addrRepo domain.AddressRepository
		logger   Logger
	}

	// Address structure
	Address struct {
		*domain.Address
	}

	// AddressID struct
	AddressID struct {
		ID string `json:"id"`
	}

	// AddressCreate struct
	AddressCreate struct {
		Line1   string `json:"line_1"`
		Line2   string `json:"line_2"`
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
		ZipCode string `json:"zip_code"`
	}

	// AddressUpdate struct
	AddressUpdate struct {
		AddressID
		AddressCreate
	}
)

// NewAddressInteractor is an AddressInteractor factory
func NewAddressInteractor(addrRepo domain.AddressRepository) *AddressInteractor {
	return &AddressInteractor{
		addrRepo: addrRepo,
	}
}

// Wrapper for domain.Address
func (i *AddressInteractor) wrap(addr *domain.Address) *Address {
	return &Address{addr}
}

// GetByID is a use case handler
func (i *AddressInteractor) GetByID(r *AddressID) (*Address, error) {
	addr, err := i.addrRepo.GetByID(r.ID)
	if err != nil {
		return nil, err
	}
	return i.wrap(addr), nil
}

// Create is a use case handler
func (i *AddressInteractor) Create(r *AddressCreate) (*Address, error) {
	addr := domain.NewAddress(r.Line1, r.Line2, r.City, r.State, r.Country, r.ZipCode)
	if err := i.addrRepo.Store(addr); err != nil {
		return nil, err
	}
	return i.wrap(addr), nil
}

// Update is a use case handler
func (i *AddressInteractor) Update(r *AddressUpdate) (*Address, error) {
	addr, err := i.addrRepo.GetByID(r.ID)
	if err != nil {
		return nil, err
	}

	addr.Line1 = r.Line1
	addr.Line2 = r.Line2
	addr.City = r.City
	addr.State = r.State
	addr.Country = r.Country
	addr.ZipCode = r.ZipCode

	if err := i.addrRepo.Update(addr); err != nil {
		return nil, err
	}
	return i.wrap(addr), nil
}

// Delete is a use case handler
func (i *AddressInteractor) Delete(r *AddressID) error {
	return i.addrRepo.Delete(r.ID)
}

// ForceDelete is a use case handler
func (i *AddressInteractor) ForceDelete(r *AddressID) error {
	return i.addrRepo.ForceDelete(r.ID)
}

// ValidationRules func returns the map with validation rules
func (r *AddressID) ValidationRules() map[string][]string {
	return map[string][]string{
		"id": []string{"required"},
	}
}

// ValidationMessages func returns map with custom validation messages
func (r *AddressID) ValidationMessages() map[string][]string {
	return nil
}

// ValidationRules func returns the map with validation rules
func (r *AddressCreate) ValidationRules() map[string][]string {
	return map[string][]string{
		"line_1":   []string{"required", "min:5", "max:250"},
		"city":     []string{"required", "min:2", "max:100"},
		"country":  []string{"required", "min:3", "max:100"},
		"zip_code": []string{"required", "min:5", "max:10"},
	}
}

// ValidationMessages func returns map with custom validation messages
func (r *AddressCreate) ValidationMessages() map[string][]string {
	return nil
}
