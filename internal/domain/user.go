package domain

import "time"

type User struct {
	ID           int     `json:"id" db:"id"`
	Username     string  `json:"username" db:"username"`
	Email        string  `json:"email" db:"email"`
	HashPassword string  `json:"-" db:"hash_password"`
	FirstName    string  `json:"first_name" db:"first_name"`
	Surname      string  `json:"surname" db:"surname"`
	Avatar       *string `json:"avatar" db:"avatar"`
	Roles        []Role  `json:"roles,omitempty" db:"-"`

	Address []Address `json:"address,omitempty" db:"-"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
