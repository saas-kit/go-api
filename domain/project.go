package domain

// Predefined user roles in project
const (
	RoleOwner  string = "owner"
	RoleAdmin  string = "admin"
	RoleRWUser string = "rw"
	RoleROUser string = "ro"
)

type (
	// ProjectRepository interface
	ProjectRepository interface {
		GetByID(id string) (Project, error)
		GetByOwnerID(id string) ([]Project, error)
		GetByMemberID(id string) ([]Project, error)
		GetAll() ([]Project, error)

		Store(*Project) error
		Update(*Project) error
		Patch(id string, data map[string]interface{}) error
		Delete(id string) error

		AddMember(projectID, memberID string, role Role) error
		DisableMember(projectID, memberID string) error
		RemoveMember(projectID, memberID string) error
	}

	// Project struct
	Project struct {
		ID       string
		Title    string
		Disabled bool
		Owner    User
		Members  []Member
	}

	// Member structure
	Member struct {
		User
		Role Role `json:"role"`
	}

	// Role type
	Role string
)
