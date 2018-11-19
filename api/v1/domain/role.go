package domain

// Role type
type Role string

// Predefined user roles in project
const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleRWUser Role = "rw"
	RoleROUser Role = "ro"
)

// IsValid determines if role is valid
func (r Role) IsValid() bool {
	empty := struct{}{}
	roles := map[Role]struct{}{
		RoleOwner:  empty,
		RoleAdmin:  empty,
		RoleRWUser: empty,
		RoleROUser: empty,
	}
	if _, ok := roles[r]; ok {
		return true
	}
	return false
}

// String converts role to string
func (r Role) String() string {
	return string(r)
}
