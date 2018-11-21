package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"saas-kit-api/api/v1/domain"
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
func (r *UserRepository) GetByID(id string) (domain.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND deleted_at is NULL LIMIT 1;", usersTableName)
	user := domain.User{}
	if err := r.db.Get(&user, query, id); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// GetByEmail retrieve user by email
func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = ? AND deleted_at is NULL LIMIT 1;", usersTableName)
	user := domain.User{}
	if err := r.db.Get(&user, query, email); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// GetList of users
func (r *UserRepository) GetList(limit, offset int) ([]domain.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at is NULL LIMIT ?, ? ;", usersTableName)
	users := make([]domain.User, 0)
	if err := r.db.Select(&users, query, offset, limit); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return users, nil
}

// GetUsersListByIDs returns users list by given ids
func (r *UserRepository) GetUsersListByIDs(ids []string) ([]domain.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (?) deleted_at is NULL;", usersTableName)
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0)
	if err := r.db.Select(&users, query, args...); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return users, nil
}

// Store new user
func (r *UserRepository) Store(user *domain.User) error {
	query := fmt.Sprintf("INSERT INTO %s (`id`, `first_name`, `last_name`, `email`, `password`, `salt`, `created_at`, `reset_token_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", usersTableName)
	if _, err := r.db.Exec(query, user.ID, user.FirstName, user.LastName,
		user.Email, user.Password, user.Salt, time.Now().Unix(), time.Now().Unix()); err != nil {
		return err
	}
	return nil
}

// Update user
func (r *UserRepository) Update(user *domain.User) error {
	query := fmt.Sprintf("UPDATE %s SET `first_name`=?, `last_name`=?, `email`=?, `password`=?, `salt`=?, `updated_at`=? WHERE id=?;", usersTableName)
	if _, err := r.db.Exec(query, user.FirstName, user.LastName,
		user.Email, user.Password, user.Salt, time.Now().Unix(), user.ID); err != nil {
		return err
	}
	return nil
}

// UpdateResetTokenTime helper
func (r *UserRepository) UpdateResetTokenTime(id string) error {
	query := fmt.Sprintf("UPDATE %s SET `reset_token_at`=?, `updated_at`=? WHERE id=?;", usersTableName)
	if _, err := r.db.Exec(query, time.Now().Unix(), time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}

// Delete record
func (r *UserRepository) Delete(id string) error {
	query := fmt.Sprintf("UPDATE %s SET `deleted_at`=? WHERE id=?;", usersTableName)
	if _, err := r.db.Exec(query, time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}
