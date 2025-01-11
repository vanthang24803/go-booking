package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type PhotoRepository interface {
	Insert(req *domain.Photo) (*domain.Photo, error)
	FindAllForListing(listingID int) ([]*domain.Photo, error)
	Remove(id string) error
}

type photoRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewPhotoRepository(db *sqlx.DB, ctx context.Context) *photoRepository {
	return &photoRepository{db: db, ctx: ctx}
}

func (r *photoRepository) FindAllForListing(listingID int) ([]*domain.Photo, error) {
	query := `
			SELECT id, listing_id, public_id, url, created_at
			FROM photos
			WHERE listing_id = $1
		`

	var photos []*domain.Photo

	err := r.db.SelectContext(r.ctx, &photos, query, listingID)

	if err != nil {
		return nil, fmt.Errorf("error finding photos for listing: %w", err)
	}

	return photos, nil
}

func (r *photoRepository) Insert(photo *domain.Photo) (*domain.Photo, error) {
	query := `
			INSERT INTO photos (listing_id, public_id, url, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at	
		`

	now := time.Now()

	err := r.db.QueryRowxContext(r.ctx, query,
		photo.ListingID,
		photo.PublicID,
		photo.URL,
		now,
	).Scan(&photo.ID, &photo.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting photo: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) Remove(id string) error {
	query := `
			DELETE FROM photos
			WHERE public_id = $1
		`

	_, err := r.db.ExecContext(r.ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting photo: %w", err)
	}

	return nil
}
