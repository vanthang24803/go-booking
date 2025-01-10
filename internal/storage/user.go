package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type UserStorage interface {
	Insert(user *domain.User) (*domain.User, error)
	Remove(id int) error
	FindOneById(id int) (*domain.User, error)
	FindOneByEmail(email string) (*domain.User, error)
	FindOneByUsername(username string) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
}

type UserRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewUserRepository(db *sqlx.DB, ctx context.Context) *UserRepository {
	return &UserRepository{db: db, ctx: ctx}
}

func (r *UserRepository) Insert(user *domain.User) (*domain.User, error) {
	query := `
        INSERT INTO users (username, email, hash_password, first_name, surname, avatar, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowxContext(r.ctx, query,
		user.Username,
		user.Email,
		user.HashPassword,
		user.FirstName,
		user.Surname,
		user.Avatar,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users
		SET username = $1, email = $2, hash_password = $3, first_name = $4, surname = $5, avatar = $6, updated_at = $7, email_verify = $8  
		WHERE id = $9
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	user.UpdatedAt = now

	err := r.db.QueryRowxContext(r.ctx, query,
		user.Username,
		user.Email,
		user.HashPassword,
		user.FirstName,
		user.Surname,
		user.Avatar,
		user.UpdatedAt,
		user.EmailVerify,
		user.ID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Remove(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(r.ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindOneByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email_verify, email, hash_password, first_name, surname, avatar, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRowxContext(r.ctx, query, email).StructScan(&user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindOneByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, hash_password, first_name, surname, avatar, created_at, updated_at , email_verify
		FROM users
		WHERE username = $1
	`

	var user domain.User
	err := r.db.QueryRowxContext(r.ctx, query, username).StructScan(&user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindOneById(id int) (*domain.User, error) {
	query := `
        SELECT id, username, email, hash_password, first_name, surname, avatar, created_at, updated_at , email_verify
        FROM users
        WHERE id = $1
    `

	var user domain.User
	err := r.db.QueryRowxContext(r.ctx, query, id).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}

	return &user, nil
}
