package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type CatalogRepository interface {
	FindAll(page int, limit int) ([]*domain.Catalog, int, error)
	FindById(id int) (*domain.Catalog, error)
	FindByName(name string) (*domain.Catalog, error)
	Insert(catalog *domain.Catalog) (*domain.Catalog, error)
	Update(catalog *domain.Catalog) (*domain.Catalog, error)
	Remove(id int) error
}

type catalogRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewCatalogRepository(db *sqlx.DB, ctx context.Context) *catalogRepository {
	return &catalogRepository{db: db, ctx: ctx}
}

func (r *catalogRepository) FindAll(page int, limit int) ([]*domain.Catalog, int, error) {
	var catalogs []*domain.Catalog
	var total int

	offset := (page - 1) * limit

	query := "SELECT id, name FROM catalogs  ORDER BY id ASC LIMIT $1 OFFSET $2"
	err := r.db.SelectContext(r.ctx, &catalogs, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalQuery := "SELECT COUNT(*) FROM catalogs"
	err = r.db.GetContext(r.ctx, &total, totalQuery)
	if err != nil {
		return nil, 0, err
	}

	return catalogs, total, nil
}

func (r *catalogRepository) FindById(id int) (*domain.Catalog, error) {
	var catalog domain.Catalog

	query := "SELECT id, name FROM catalogs WHERE id = $1"

	err := r.db.GetContext(r.ctx, &catalog, query, id)
	if err != nil {
		return nil, err
	}

	return &catalog, nil
}

func (r *catalogRepository) FindByName(name string) (*domain.Catalog, error) {
	var catalog domain.Catalog

	query := "SELECT id, name FROM catalogs WHERE name = $1"

	err := r.db.GetContext(r.ctx, &catalog, query, name)

	if err != nil {
		return nil, err
	}

	return &catalog, nil
}

func (r *catalogRepository) Insert(catalog *domain.Catalog) (*domain.Catalog, error) {
	query := "INSERT INTO catalogs (name) VALUES ($1) RETURNING id, name"

	err := r.db.GetContext(r.ctx, catalog, query, catalog.Name)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (r *catalogRepository) Update(catalog *domain.Catalog) (*domain.Catalog, error) {
	query := "UPDATE catalogs SET name = $1 WHERE id = $2 RETURNING id, name"

	err := r.db.GetContext(r.ctx, catalog, query, catalog.Name, catalog.ID)

	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (r *catalogRepository) Remove(id int) error {
	query := "DELETE FROM catalogs WHERE id = $1"

	_, err := r.db.ExecContext(r.ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}