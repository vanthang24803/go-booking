package storage

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRoleStorage_FindRoleByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRoleRepository(sqlxDB, context.Background())

	mockRole := domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `SELECT id, role_name, description, created_at, updated_at FROM roles WHERE role_name = \$1`
	rows := sqlmock.NewRows([]string{"id", "role_name", "description", "created_at", "updated_at"}).
		AddRow(mockRole.ID, mockRole.RoleName, mockRole.Description, mockRole.CreatedAt, mockRole.UpdatedAt)

	mock.ExpectQuery(query).WithArgs("admin").WillReturnRows(rows)

	role, err := repo.FindRoleByName("admin")
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "admin", role.RoleName)
}

func TestRoleStorage_InsertRoleToUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRoleRepository(sqlxDB, context.Background())

	userRole := &domain.UserRole{
		UserID: 1,
		RoleID: 2,
	}

	query := `INSERT INTO user_roles \(user_id, role_id\) VALUES \(\$1, \$2\)`
	mock.ExpectExec(query).WithArgs(userRole.UserID, userRole.RoleID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertRoleToUser(userRole)
	assert.NoError(t, err)
}

func TestRoleStorage_FindRolesByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRoleRepository(sqlxDB, context.Background())

	mockRoles := []domain.Role{
		{ID: 1, RoleName: "admin", Description: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, RoleName: "user", Description: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	query := `
		SELECT r.id, r.role_name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = \$1
	`
	rows := sqlmock.NewRows([]string{"id", "role_name", "description", "created_at", "updated_at"})
	for _, role := range mockRoles {
		rows.AddRow(role.ID, role.RoleName, role.Description, role.CreatedAt, role.UpdatedAt)
	}

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	roles, err := repo.FindRolesByUser(1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(roles))
	assert.Equal(t, "admin", roles[0].RoleName)
}

func TestRoleStorage_FindRoleById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRoleRepository(sqlxDB, context.Background())

	mockRole := domain.Role{
		ID:          1,
		RoleName:    "admin",
		Description: nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `SELECT id, role_name, description, created_at, updated_at FROM roles WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"id", "role_name", "description", "created_at", "updated_at"}).
		AddRow(mockRole.ID, mockRole.RoleName, mockRole.Description, mockRole.CreatedAt, mockRole.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	role, err := repo.FindRoleById(1)
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "admin", role.RoleName)
}

func TestRoleStorage_RemoveRoleFromUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRoleRepository(sqlxDB, context.Background())

	userRole := &domain.UserRole{
		UserID: 1,
		RoleID: 2,
	}

	query := `DELETE FROM user_roles WHERE user_id = \$1 AND role_id = \$2`
	mock.ExpectExec(query).WithArgs(userRole.UserID, userRole.RoleID).WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.RemoveRoleFromUser(userRole)
	assert.NoError(t, err)
}
