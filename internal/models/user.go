package models

import "time"

type User struct {
	ID          string       `json:"id"`
	FullName    string       `json:"full_name"`
	Email       string       `json:"email"`
	Password    string       `json:"-"`
	RoleID      string       `json:"role_id"`
	Role        *Role        `json:"role,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	LastLoginAt *time.Time   `json:"last_login_at,omitempty"`
}

// HasPermission checks whether the user has the given resource:action permission.
func (u *User) HasPermission(resource, action string) bool {
	for _, p := range u.Permissions {
		if p.Resource == resource && p.Action == action {
			return true
		}
	}
	return false
}
