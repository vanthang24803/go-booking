package domain

import "time"

type Address struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"-" db:"user_id"`
	Street     string `json:"street" db:"street"`
	City       string `json:"city" db:"city"`
	State      string `json:"state" db:"state"`
	Country    string `json:"country" db:"country"`
	PostalCode string `json:"postal_code" db:"postal_code"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
