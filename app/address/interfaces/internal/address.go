package internal

import (
	"saas-kit-api/app/address/usecases"
)

type (
	// AddressService structure
	AddressService struct {
		addrIntr usecases.AddressInteractor
	}

	// Address structure
	Address struct {
		*usecases.Address
	}
)

// NewAddressService is an AddressService factory
func NewAddressService(addrIntr usecases.AddressInteractor) *AddressService {
	return &AddressService{
		addrIntr: addrIntr,
	}
}

// Wrapper for usecases.Address
func (s *AddressService) wrap(addr *usecases.Address) *Address {
	return &Address{addr}
}

// GetByID handler
func (s *AddressService) GetByID(id string) (*Address, error) {
	addr, err := s.addrIntr.GetByID(&usecases.AddressID{ID: id})
	if err != nil {
		return nil, err
	}
	return s.wrap(addr), nil
}

// Delete handler
func (s *AddressService) Delete(id string) error {
	return s.addrIntr.Delete(&usecases.AddressID{ID: id})
}
