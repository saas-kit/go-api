package usecases

import "saas-kit-api/api/v1/domain"

type (
	// Project structure
	Project struct {
		domain.Project
	}

	// ProjectInteractor structure
	ProjectInteractor struct {
		prRepo domain.ProjectRepository
	}
)

// NewProjectInteractor is factory function,
// returns a new instance of the ProjectInteractor
func NewProjectInteractor(prRepo domain.ProjectRepository) *ProjectInteractor {
	return &ProjectInteractor{
		prRepo: prRepo,
	}
}

// AddMember to a project
func (i *ProjectInteractor) AddMember(projectID, userID string, role domain.Role) error {
	if err := i.prRepo.AddMember(projectID, userID, role); err != nil {
		return err
	}
	return nil
}
