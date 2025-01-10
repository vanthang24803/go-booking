package domain

import "time"

type Token struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Name      string     `db:"name"`
	Token     string     `db:"token"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	ExpiredAt *time.Time `db:"expired_at"`
}
