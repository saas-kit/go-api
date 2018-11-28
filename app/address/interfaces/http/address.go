package jsonrpc

import (
	"net/http"
	"saas-kit-api/app/address/usecases"
	"saas-kit-api/pkg/server"

	"github.com/labstack/echo"
)

type (
	// AddressService structure
	AddressService struct {
		addrIntr usecases.AddressInteractor
	}

	// Address structure
	Address struct {
		*usecases.Address
		Object string `json:"object,omitempty"`
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
	return &Address{addr, "address"}
}

// GetByID handler
func (s *AddressService) GetByID(c echo.Context) error {
	id := c.Param("id")
	addr, err := s.addrIntr.GetByID(&usecases.AddressID{ID: id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, server.Response(s.wrap(addr)))
}

// Create handler
func (s *AddressService) Create(c echo.Context) error {
	req := &usecases.AddressCreate{}
	if err := c.Bind(req); err != nil {
		return err
	}
	addr, err := s.addrIntr.Create(req)
	if err != nil {
		return server.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, server.Response(s.wrap(addr)))
}

// Update handler
func (s *AddressService) Update(c echo.Context) error {
	req := &usecases.AddressUpdate{}
	if err := c.Bind(req); err != nil {
		return err
	}
	addr, err := s.addrIntr.Update(req)
	if err != nil {
		return server.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, server.Response(s.wrap(addr)))
}

// Delete handler
func (s *AddressService) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := s.addrIntr.Delete(&usecases.AddressID{ID: id}); err != nil {
		return server.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, server.Response(true))
}
