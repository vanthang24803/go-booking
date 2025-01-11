package domain

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	HashPassword string    `json:"-" db:"hash_password"`
	FirstName    string    `json:"first_name" db:"first_name"`
	Surname      string    `json:"surname" db:"surname"`
	Avatar       *string   `json:"avatar" db:"avatar"`
	EmailVerify  bool      `json:"is_verify" db:"email_verify"`
	Roles        []Role    `json:"roles,omitempty" db:"-"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

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

type Token struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Name      string     `db:"name"`
	Token     string     `db:"token"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	ExpiredAt *time.Time `db:"expired_at"`
}
