package repositories

import (
	"saas-kit-api/app/auth/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	usersTableName string = "users"
)

type (
	// UserRepository struct
	UserRepository struct {
		db *sqlx.DB
	}
)

// NewUserRepository is a factory function,
// returns a new instance of the UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

// GetByID retrieve user by id
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query := "SELECT * FROM ? WHERE id = ? LIMIT 1;"
	user := &domain.User{}
	if err := r.db.Get(user, query, usersTableName, id); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieve user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := "SELECT * FROM ? WHERE email = ? LIMIT 1;"
	user := &domain.User{}
	if err := r.db.Get(user, query, usersTableName, email); err != nil {
		return nil, err
	}
	return user, nil
}

// Store new user
func (r *UserRepository) Store(user *domain.User) error {
	query := "INSERT INTO ? (`id`, `email`, `password`, `salt`, `confirmed`, `disabled`, `reset_token_at`, `created_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
	_, err := r.db.Exec(
		query, usersTableName,
		user.ID, user.Email, user.Password, user.Salt,
		user.Confirmed, user.Disabled,
		time.Now().Unix(), time.Now().Unix(),
	)
	return err
}

// ChangeEmail helper
func (r *UserRepository) ChangeEmail(id, email string) error {
	query := "UPDATE ? SET `email`=?, `confirmed`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, email, false, time.Now().Unix(), id)
	return err
}

// UpdateResetTokenTime helper
func (r *UserRepository) UpdateResetTokenTime(id string) error {
	query := "UPDATE ? SET `reset_token_at`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, time.Now().Unix(), time.Now().Unix(), id)
	return err
}

// UpdatePassword helper
func (r *UserRepository) UpdatePassword(id, password string) error {
	query := "UPDATE ? SET `password`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, password, time.Now().Unix(), id)
	return err
}

// Confirm email helper
func (r *UserRepository) Confirm(id string) error {
	query := "UPDATE ? SET `confirmed`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, true, time.Now().Unix(), id)
	return err
}

// Unconfirm email helper
func (r *UserRepository) Unconfirm(id string) error {
	query := "UPDATE ? SET `confirmed`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, false, time.Now().Unix(), id)
	return err
}

// Disable user helper
func (r *UserRepository) Disable(id string) error {
	query := "UPDATE ? SET `disabled`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, true, time.Now().Unix(), id)
	return err
}

// Enable user helper
func (r *UserRepository) Enable(id string) error {
	query := "UPDATE ? SET `disabled`=?, `updated_at`=? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, false, time.Now().Unix(), id)
	return err
}

// Delete record
func (r *UserRepository) Delete(id string) error {
	query := "DELETE FROM ? WHERE id=?;"
	_, err := r.db.Exec(query, usersTableName, id)
	return err
}
