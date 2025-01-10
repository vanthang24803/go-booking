package storage

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTokenStorage_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewTokenRepository(sqlxDB, ctx)

	token := &domain.Token{
		UserID:    1,
		Token:     "abcd1234",
		Name:      "auth_token",
		ExpiredAt: nil,
	}

	query := `
		INSERT INTO tokens \(user_id, token, name, created_at, updated_at, expired_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)
		RETURNING id, created_at, updated_at
	`
	mock.ExpectQuery(query).
		WithArgs(token.UserID, token.Token, token.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), token.ExpiredAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	insertedToken, err := repo.Insert(token)
	assert.NoError(t, err)
	assert.Equal(t, 1, insertedToken.ID)
	assert.NotNil(t, insertedToken.CreatedAt)
	assert.NotNil(t, insertedToken.UpdatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenStorage_Remove(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewTokenRepository(sqlxDB, ctx)

	query := `DELETE FROM tokens WHERE id = \$1`
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Remove(1)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenStorage_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewTokenRepository(sqlxDB, ctx)

	now := time.Now()

	token := &domain.Token{
		ID:        1,
		UserID:    1,
		Token:     "updated_token_123",
		Name:      "updated_auth_token",
		CreatedAt: now,
		UpdatedAt: now,
		ExpiredAt: nil,
	}

	query := regexp.QuoteMeta(`
        UPDATE tokens
        SET user_id = $1, token = $2, created_at = $3, updated_at = $4, expired_at = $5, name= $6
        WHERE id = $7
        RETURNING id, created_at, updated_at
    `)

	mock.ExpectQuery(query).
		WithArgs(token.UserID, token.Token, token.CreatedAt, sqlmock.AnyArg(), token.ExpiredAt, token.Name, token.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(token.ID, token.CreatedAt, time.Now()))

	updatedToken, err := repo.Update(token)
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("error was not expected while updating token: %s", err)
	}

	assert.Equal(t, token.ID, updatedToken.ID)
	assert.Equal(t, token.UserID, updatedToken.UserID)
	assert.Equal(t, token.Token, updatedToken.Token)
	assert.Equal(t, token.Name, updatedToken.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenStorage_FindOneByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewTokenRepository(sqlxDB, ctx)

	tokenName := "test_token"
	userId := 1
	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	expectedToken := &domain.Token{
		ID:        1,
		UserID:    userId,
		Token:     "token_value_123",
		Name:      tokenName,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiredAt: &expiredAt,
	}

	query := regexp.QuoteMeta(`
        SELECT id, user_id, token, name, created_at, updated_at, expired_at
        FROM tokens
        WHERE name = $1 AND user_id = $2
    `)

	rows := sqlmock.NewRows([]string{"id", "user_id", "token", "name", "created_at", "updated_at", "expired_at"}).
		AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.Token, expectedToken.Name, expectedToken.CreatedAt, expectedToken.UpdatedAt, expectedToken.ExpiredAt)

	mock.ExpectQuery(query).WithArgs(tokenName, userId).WillReturnRows(rows)

	token, err := repo.FindOneByToken(tokenName, userId)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, expectedToken.ID, token.ID)
	assert.Equal(t, expectedToken.UserID, token.UserID)
	assert.Equal(t, expectedToken.Token, token.Token)
	assert.Equal(t, expectedToken.Name, token.Name)
	assert.Equal(t, expectedToken.CreatedAt, token.CreatedAt)
	assert.Equal(t, expectedToken.UpdatedAt, token.UpdatedAt)
	assert.Equal(t, expectedToken.ExpiredAt, token.ExpiredAt)

	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery(query).WithArgs(tokenName, userId).WillReturnError(sql.ErrNoRows)

	token, err = repo.FindOneByToken(tokenName, userId)
	assert.NoError(t, err)
	assert.Nil(t, token)

	assert.NoError(t, mock.ExpectationsWereMet())
}
