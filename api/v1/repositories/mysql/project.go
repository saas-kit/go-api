package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"saas-kit-api/api/v1/domain"
)

const (
	projectTableName string = "projects"
	membersTableName string = "project_members"
)

// Predefined errors
var (
	ErrMissedOwner         = errors.New("Missed project owner")
	ErrMemberAlreadyExists = errors.New("Member already exists")
	ErrCouldNotChangeOwner = errors.New("Could not change project owner")
	ErrUserShouldBeMember  = errors.New("The user is not a member of the project")
	ErrUserAlreadyOwner    = errors.New("The user is already owner of the project")
)

type (
	// ProjectRepository struct
	ProjectRepository struct {
		db       *sqlx.DB
		userRepo userRepository
	}

	// user repository interface for injection to the ProjectRepository
	userRepository interface {
		GetUsersListByIDs([]string) ([]domain.User, error)
	}

	// internal member structure
	member struct {
		ProjectID string      `json:"project_id"`
		UserID    string      `json:"user_id"`
		Role      domain.Role `json:"role"`
		Disabled  bool        `json:"disabled"`
	}
)

// NewProjectRepository is a factory function,
// returns a new instance of the ProjectRepository
func NewProjectRepository(db *sqlx.DB, ur userRepository) *ProjectRepository {
	return &ProjectRepository{db, ur}
}

// GetByID retrieve project by id
func (r *ProjectRepository) GetByID(id string) (*domain.Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? LIMIT 1;", projectTableName)
	project := domain.Project{}
	if err := r.db.Get(&project, query, id); err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByOwnerID retrieve project by owner id
func (r *ProjectRepository) GetByOwnerID(ownerID string) ([]domain.Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (SELECT project_id FROM %s WHERE role = ? AND user_id = ? AND disabled = 0);", projectTableName, membersTableName)
	projects := make([]domain.Project, 0)
	if err := r.db.Select(&projects, query, domain.RoleOwner.String(), ownerID); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return projects, nil
}

// GetByMemberID retrieve project by member id
func (r *ProjectRepository) GetByMemberID(memberID string) ([]domain.Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id IN (SELECT project_id FROM %s WHERE user_id = ? AND disabled = 0);", projectTableName, membersTableName)
	projects := make([]domain.Project, 0)
	if err := r.db.Select(&projects, query, memberID); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return projects, nil
}

// GetList of projects
func (r *ProjectRepository) GetList(limit, offset int) ([]domain.Project, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ?, ?;", projectTableName)
	projects := make([]domain.Project, 0)
	if err := r.db.Select(&projects, query, offset, limit); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return projects, nil
}

// Store new project
func (r *ProjectRepository) Store(project *domain.Project) error {
	if project.Owner.ID == "" {
		return ErrMissedOwner
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (`id`, `title`, `disabled`, `created_at`) VALUES (?, ?, ?, ?);", projectTableName)
	if _, err := tx.Exec(query, project.ID, project.Title, project.Disabled, time.Now().Unix()); err != nil {
		tx.Rollback()
		return err
	}
	query2 := fmt.Sprintf("INSERT INTO %s (`project_id`, `user_id`, `role`) VALUES (?, ?, ?);", membersTableName)
	if _, err := tx.Exec(query2, project.ID, project.Owner.ID, domain.RoleOwner.String()); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Update project
func (r *ProjectRepository) Update(project *domain.Project) error {
	query := fmt.Sprintf("UPDATE %s SET `title`=?, `disabled`=?, `updated_at`=? WHERE id=?;", projectTableName)
	if _, err := r.db.Exec(query, project.Title, project.Disabled, time.Now().Unix(), project.ID); err != nil {
		return err
	}
	return nil
}

// Delete record
func (r *ProjectRepository) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?;", projectTableName)
	if _, err := r.db.Exec(query, time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}

// AddMember to the project
func (r *ProjectRepository) AddMember(projectID, memberID string, role domain.Role) error {
	query := fmt.Sprintf("SELECT count(project_id) FROM %s WHERE project_id = ? AND user_id = ?;", membersTableName)
	var count int
	if err := r.db.Get(&count, query, projectID, memberID); err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		return ErrMemberAlreadyExists
	}
	query2 := fmt.Sprintf("INSERT INTO %s (`project_id`, `user_id`, `role`) VALUES (?, ?, ?);", membersTableName)
	if _, err := r.db.Exec(query2, projectID, memberID, role.String()); err != nil {
		return err
	}
	return nil
}

// GetMembers list by project id
func (r *ProjectRepository) GetMembers(projectID string) ([]domain.Member, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE project_id = ?;", membersTableName)
	members := make([]member, 0)
	if err := r.db.Select(members, query, projectID); err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err != sql.ErrNoRows {
		return nil, nil
	}
	ids := getMembersIDs(members)
	if len(ids) == 0 {
		return nil, nil
	}
	users, err := r.userRepo.GetUsersListByIDs(ids)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	muMap := make(map[string]domain.Role, 0)
	for _, m := range members {
		muMap[m.UserID] = m.Role
	}
	return domain.GenMembers(users, muMap), nil
}

// DisableMember of project
func (r *ProjectRepository) DisableMember(projectID, memberID string) error {
	return r.toggleMember(projectID, memberID, true)
}

// EnableMember of project
func (r *ProjectRepository) EnableMember(projectID, memberID string) error {
	return r.toggleMember(projectID, memberID, false)
}

// Toggle membership
func (r *ProjectRepository) toggleMember(projectID, memberID string, disabled bool) error {
	query := fmt.Sprintf("UPDATE %s SET `disabled` = ? WHERE project_id = ? AND user_id = ?;", membersTableName)
	if _, err := r.db.Exec(query, disabled, projectID, memberID); err != nil {
		return err
	}
	return nil
}

// ChangeMemberRole in a project
func (r *ProjectRepository) ChangeMemberRole(projectID, memberID string, role domain.Role) error {
	if role == domain.RoleOwner {
		return ErrCouldNotChangeOwner
	}
	query := fmt.Sprintf("UPDATE %s SET `role` = ? WHERE project_id = ? AND user_id = ?;", membersTableName)
	if _, err := r.db.Exec(query, role.String(), projectID, memberID); err != nil {
		return err
	}
	return nil
}

// TransferProject to another user
func (r *ProjectRepository) TransferProject(projectID, memberID string) error {
	query := fmt.Sprintf("SELECT role FROM %s WHERE project_id = ? AND user_id = ? LIMIT 1;", membersTableName)
	var role domain.Role
	if err := r.db.Get(&role, query, projectID, memberID); err != nil && err != sql.ErrNoRows {
		return err
	} else if err == sql.ErrNoRows {
		return ErrUserShouldBeMember
	}
	if role == domain.RoleOwner {
		return ErrUserAlreadyOwner
	}
	tx := r.db.MustBegin()
	query2 := fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND role = ?;", membersTableName)
	if _, err := tx.Exec(query2, projectID, domain.RoleOwner.String()); err != nil {
		tx.Rollback()
		return err
	}
	query3 := fmt.Sprintf("UPDATE %s SET `role` = ? WHERE project_id = ? AND user_id = ?;", membersTableName)
	if _, err := tx.Exec(query3, domain.RoleOwner.String(), projectID, memberID); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// RemoveMember from a project
func (r *ProjectRepository) RemoveMember(projectID, memberID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND user_id = ?;", membersTableName)
	if _, err := r.db.Exec(query, projectID, memberID); err != nil {
		return err
	}
	return nil
}

// returns members ids
func getMembersIDs(members []member) []string {
	ids := make([]string, len(members))
	for key, val := range members {
		ids[key] = val.UserID
	}
	return ids
}
