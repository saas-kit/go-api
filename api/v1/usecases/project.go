package usecases

import "saas-kit-api/api/v1/domain"

type (
	// Project structure
	Project struct {
		*domain.Project
	}

	// Member structure
	Member struct {
		domain.Member
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

// GetByID retrieves project by id
func (i *ProjectInteractor) GetByID(id string) (*Project, error) {
	pr, err := i.prRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return i.wrap(pr), nil
}

// GetByMemberID retrieves project by project member id
func (i *ProjectInteractor) GetByMemberID(userID string) ([]Project, error) {
	prs, err := i.prRepo.GetByMemberID(userID)
	if err != nil {
		return nil, err
	}
	return i.wrapList(prs), nil
}

// GetMembers by project id
func (i *ProjectInteractor) GetMembers(id string) ([]Member, error) {
	mbrs, err := i.prRepo.GetMembers(id)
	if err != nil {
		return nil, err
	}
	return i.wrapMembersList(mbrs), nil
}

// AddMember to a project
func (i *ProjectInteractor) AddMember(projectID, userID string, role domain.Role) error {
	if err := i.prRepo.AddMember(projectID, userID, role); err != nil {
		return err
	}
	return nil
}

// wrapper for domain.Project structure
func (i *ProjectInteractor) wrap(pr *domain.Project) *Project {
	return &Project{pr}
}

// wrapper for domain.Project structures list
func (i *ProjectInteractor) wrapList(prs []domain.Project) []Project {
	list := make([]Project, 0, len(prs))
	for _, p := range prs {
		list = append(list, Project{&p})
	}
	return list
}

// wrapper for domain.Member structures list
func (i *ProjectInteractor) wrapMembersList(mbrs []domain.Member) []Member {
	list := make([]Member, 0, len(mbrs))
	for _, m := range mbrs {
		list = append(list, Member{m})
	}
	return list
}
