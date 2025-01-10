package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewUserRepository(sqlxDB, ctx)

	user := &domain.User{
		Username:     faker.Username(),
		Email:        faker.Email(),
		HashPassword: faker.Password(),
		FirstName:    faker.FirstName(),
		Surname:      faker.LastName(),
		Avatar:       nil,
	}

	query := `
        INSERT INTO users \(username, email, hash_password, first_name, surname, avatar, created_at, updated_at\)
        VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery(query).
		WithArgs(
			user.Username,
			user.Email,
			user.HashPassword,
			user.FirstName,
			user.Surname,
			user.Avatar,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(rows)

	createdUser, err := repo.Insert(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, createdUser.ID)
	assert.WithinDuration(t, now, createdUser.CreatedAt, time.Second)
	assert.WithinDuration(t, now, createdUser.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOneByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewUserRepository(sqlxDB, ctx)

	email := faker.Email()
	user := &domain.User{
		ID:           1,
		Username:     faker.Username(),
		Email:        email,
		HashPassword: faker.Password(),
		FirstName:    faker.Username(),
		Surname:      faker.LastName(),
		Avatar:       nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		SELECT id, username, email, hash_password, first_name, surname, avatar, created_at, updated_at
		FROM users
		WHERE email = \$1
	`

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "hash_password", "first_name", "surname", "avatar", "created_at", "updated_at",
	}).AddRow(
		user.ID,
		user.Username,
		user.Email,
		user.HashPassword,
		user.FirstName,
		user.Surname,
		user.Avatar,
		user.CreatedAt,
		user.UpdatedAt,
	)

	mock.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)

	foundUser, err := repo.FindOneByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewUserRepository(sqlxDB, ctx)

	userId := 1
	query := `
		DELETE FROM users
		WHERE id = \$1
	`

	mock.ExpectExec(query).WithArgs(userId).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Remove(userId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOneByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	ctx := context.Background()
	repo := NewUserRepository(sqlxDB, ctx)

	userId := 1
	query := `
        SELECT id, username, email, hash_password, first_name, surname, avatar, created_at, updated_at
        FROM users
        WHERE id = \$1
    `

	mock.ExpectQuery(query).WithArgs(userId).WillReturnError(sql.ErrNoRows)

	foundUser, err := repo.FindOneById(userId)

	assert.Nil(t, foundUser)
	assert.Error(t, err)
	assert.EqualError(t, err, "user with id 1 not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}
