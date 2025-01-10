package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type RoleStorage interface {
	InsertRoleToUser(userRole *domain.UserRole) error
	FindRolesByUser(userId int) ([]domain.Role, error)
	FindRoleByName(name string) (*domain.Role, error)
	FindRoleById(id int) (*domain.Role, error)
	RemoveRoleFromUser(userRole *domain.UserRole) error
}

type RoleRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewRoleRepository(db *sqlx.DB, ctx context.Context) *RoleRepository {
	return &RoleRepository{db: db, ctx: ctx}
}

func (r *RoleRepository) FindRoleByName(name string) (*domain.Role, error) {
	query := `
		SELECT id, role_name, description, created_at, updated_at
		FROM roles
		WHERE role_name = $1
	`
	var role domain.Role
	err := r.db.GetContext(r.ctx, &role, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return &role, nil
}

func (r *RoleRepository) InsertRoleToUser(userRole *domain.UserRole) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(r.ctx, query, userRole.UserID, userRole.RoleID)
	if err != nil {
		return fmt.Errorf("failed to insert role to user: %w", err)
	}
	return nil
}

func (r *RoleRepository) FindRolesByUser(userId int) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.role_name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`
	var roles []domain.Role
	err := r.db.SelectContext(r.ctx, &roles, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles by user: %w", err)
	}
	return roles, nil
}

func (r *RoleRepository) FindRoleById(id int) (*domain.Role, error) {
	query := `
		SELECT id, role_name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`
	var role domain.Role
	err := r.db.GetContext(r.ctx, &role, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return &role, nil
}

func (r *RoleRepository) RemoveRoleFromUser(userRole *domain.UserRole) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2
	`
	_, err := r.db.ExecContext(r.ctx, query, userRole.UserID, userRole.RoleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}
