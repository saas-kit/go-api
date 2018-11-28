package repositories

import (
	"saas-kit-api/app/address/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// AddressRepository struct
	AddressRepository struct {
		db        *sqlx.DB
		tableName string
	}
)

// NewAddressRepository factory
func NewAddressRepository(db *sqlx.DB, tableName ...string) *AddressRepository {
	repo := &AddressRepository{db, "addresses"}
	if len(tableName) > 0 {
		repo.tableName = tableName[0]
	}
	return repo
}

// GetByID retrieve record by address id
func (r *AddressRepository) GetByID(id string) (*domain.Address, error) {
	query := "SELECT * FROM ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;"
	user := &domain.Address{}
	if err := r.db.Get(user, query, r.tableName, id); err != nil {
		return nil, err
	}
	return user, nil
}

// Store new address
func (r *AddressRepository) Store(addr *domain.Address) error {
	query := "INSERT INTO ? (id, line_1, line_2, city, state, country, zip_code, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := r.db.Exec(
		query, r.tableName,
		addr.ID,
		addr.Line1, addr.Line2,
		addr.City, addr.State,
		addr.Country, addr.ZipCode,
		time.Now().Unix(),
	)
	return err
}

// Update an address record
func (r *AddressRepository) Update(addr *domain.Address) error {
	query := "UPDATE ? SET line_1=?, line_2=?, city=?, state=?, country=?, zip_code=?, updated_at=? WHERE id=?;"
	_, err := r.db.Exec(
		query, r.tableName,
		addr.Line1, addr.Line2,
		addr.City, addr.State,
		addr.Country, addr.ZipCode,
		time.Now().Unix(),
		addr.ID,
	)
	return err
}

// Delete an address record
func (r *AddressRepository) Delete(id string) error {
	query := "UPDATE ? SET deleted_at=? WHERE id=?;"
	_, err := r.db.Exec(
		query, r.tableName,
		time.Now().Unix(),
		id,
	)
	return err
}

// ForceDelete address
func (r *AddressRepository) ForceDelete(id string) error {
	query := "DELETE FROM ? WHERE id=?;"
	_, err := r.db.Exec(query, r.tableName, id)
	return err
}
