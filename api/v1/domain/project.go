package domain

type (
	// ProjectRepository interface
	ProjectRepository interface {
		GetByID(id string) (*Project, error)
		GetByOwnerID(id string) ([]Project, error)
		GetByMemberID(id string) ([]Project, error)
		GetList(limit, offset int) ([]Project, error)

		Store(*Project) error
		Update(*Project) error
		Delete(id string) error

		GetMembers(projectID string) ([]Member, error)
		AddMember(projectID, memberID string, role Role) error
		DisableMember(projectID, memberID string) error
		RemoveMember(projectID, memberID string) error
		ChangeMemberRole(projectID, memberID string, role Role) error
		TransferProject(projectID, memberID string) error
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
)

// NewMember is a factory function
func NewMember(user User, role Role) Member {
	return Member{user, role}
}

// GenMembers generate members objects from array of users
func GenMembers(users []User, userRole map[string]Role) []Member {
	members := make([]Member, 0)
	for _, u := range users {
		if role, ok := userRole[u.ID]; ok {
			members = append(members, NewMember(u, role))
		}
	}
	return members
}
