package entities

const (
	// AdministratorRole defines role of administrator
	AdministratorRole = "administrator"
)

// User defines a user for our application
type User struct {
	ID        uint     `json:"id"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FirstName string   `josn:"firstname"`
	LastName  string   `json:"lastname"`
	Roles     []string `json:"roles"`
	IsOnline  bool     `json:"isonline"`
	IsEnabled bool     `json:"isenabled"`
}

// HasRole returns true if the user is in the role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role == roleName {
			return true
		}
	}
	return false
}
