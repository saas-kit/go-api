package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"saas-kit-api/api/v1/domain"
)

const (
	invitationTableName string = "invitations"
)

type (
	// InvitationsRepository struct
	InvitationsRepository struct {
		db *sqlx.DB
	}
)

// NewInvitationsRepository is a factory function,
// returns a new instance of the InvitationsRepository
func NewInvitationsRepository(db *sqlx.DB) *InvitationsRepository {
	return &InvitationsRepository{db}
}

// GetByID retrieve invitations by id
func (r *InvitationsRepository) GetByID(id string) (*domain.Invitation, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1;", invitationTableName)
	invitation := domain.Invitation{}
	if err := r.db.Get(&invitation, query, id); err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetByEmail retrieve invitations by email
func (r *InvitationsRepository) GetByEmail(email string) ([]domain.Invitation, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = ?;", invitationTableName)
	invitations := make([]domain.Invitation, 0)
	if err := r.db.Select(&invitations, query, email); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return invitations, nil
}

// GetByProject retrieve invitations by project id
func (r *InvitationsRepository) GetByProject(projectID string) ([]domain.Invitation, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = ?;", invitationTableName)
	invitations := make([]domain.Invitation, 0)
	if err := r.db.Select(&invitations, query, projectID); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return invitations, nil
}

// GetList of invitations
func (r *InvitationsRepository) GetList(limit, offset int) ([]domain.Invitation, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ?, ? ;", invitationTableName)
	invitations := make([]domain.Invitation, 0)
	if err := r.db.Select(&invitations, query, offset, limit); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return invitations, nil
}

// Store new invitations
func (r *InvitationsRepository) Store(invitation *domain.Invitation) error {
	query := fmt.Sprintf("INSERT INTO %s (`id`, `email`, `project_id`, `role`, `sent_at`, `created_at`) VALUES (?, ?, ?, ?, ?, ?);", invitationTableName)
	if _, err := r.db.Exec(query, invitation.ID, invitation.Email, invitation.Project.ID,
		invitation.Role.String(), invitation.SentAt, time.Now().Unix()); err != nil {
		return err
	}
	return nil
}

// Update invitations
func (r *InvitationsRepository) Update(invitation *domain.Invitation) error {
	query := fmt.Sprintf("UPDATE %s SET `role`=?, `confirmed`=?, `sent_at`=?, `updated_at`=? WHERE id=?;", invitationTableName)
	if _, err := r.db.Exec(query, invitation.Role.String(), invitation.Confirmed, invitation.SentAt, time.Now().Unix(), invitation.ID); err != nil {
		return err
	}
	return nil
}

// Delete record
func (r *InvitationsRepository) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?;", invitationTableName)
	if _, err := r.db.Exec(query, time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}

// Confirm invitations
func (r *InvitationsRepository) Confirm(id string) error {
	query := fmt.Sprintf("UPDATE %s SET `confirmed`=?, `updated_at`=? WHERE id=?;", invitationTableName)
	if _, err := r.db.Exec(query, true, time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}
