package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type TokenStorage interface {
	Insert(token *domain.Token) (*domain.Token, error)
	Update(token *domain.Token) (*domain.Token, error)
	Remove(id int) error
	FindOneByToken(token string, userId int) (*domain.Token, error)
}

type TokenRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewTokenRepository(db *sqlx.DB, ctx context.Context) *TokenRepository {
	return &TokenRepository{db: db, ctx: ctx}
}

func (r *TokenRepository) FindOneByToken(tokenName string, userId int) (*domain.Token, error) {
	query := `
        SELECT id, user_id, token, name, created_at, updated_at, expired_at
        FROM tokens
        WHERE name = $1 AND user_id = $2
    `

	var t domain.Token
	err := r.db.QueryRowContext(r.ctx, query, tokenName, userId).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.Name,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.ExpiredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying token: %w", err)
	}

	return &t, nil
}
func (r *TokenRepository) Insert(token *domain.Token) (*domain.Token, error) {
	query := `
		INSERT INTO tokens (user_id, token, name, created_at, updated_at, expired_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	token.CreatedAt = now
	token.UpdatedAt = now

	err := r.db.QueryRowContext(r.ctx, query,
		token.UserID,
		token.Token,
		token.Name,
		token.CreatedAt,
		token.UpdatedAt,
		token.ExpiredAt,
	).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting token: %w", err)
	}

	return token, nil
}

func (r *TokenRepository) Update(token *domain.Token) (*domain.Token, error) {
	query := `
		UPDATE tokens
		SET user_id = $1, token = $2, created_at = $3, updated_at = $4, expired_at = $5, name= $6
		WHERE id = $7
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(r.ctx, query,
		token.UserID,
		token.Token,
		token.CreatedAt,
		token.UpdatedAt,
		token.ExpiredAt,
		token.Name,
		token.ID,
	).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating token: %w", err)
	}

	return token, nil
}

func (r *TokenRepository) Remove(id int) error {
	query := `
		DELETE FROM tokens
		WHERE id = $1
	`

	_, err := r.db.ExecContext(r.ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}

	return nil
}
