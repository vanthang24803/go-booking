package domain

import "time"

type Role struct {
	ID          int       `json:"id" db:"id"`
	RoleName    string    `json:"name" db:"role_name"`
	Description *string   `json:"-" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

type UserRole struct {
	UserID int `db:"user_id"`
	RoleID int `db:"role_id"`
}
