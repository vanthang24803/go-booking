package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type ListingRepository interface {
	FindAll(page int, limit int) ([]*domain.Listing, int, int, error)
	FindOne(id int) (*domain.Listing, error)
	Save(listing *domain.Listing) (*domain.Listing, error)
	Update(id int, listing *domain.Listing) (*domain.Listing, error)
	Remove(id int) error
}

type listingRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewListingRepository(db *sqlx.DB, ctx context.Context) *listingRepository {
	return &listingRepository{db: db, ctx: ctx}
}

func (r *listingRepository) FindAll(page int, limit int) ([]*domain.Listing, int, int, error) {

	var listings []*domain.Listing
	var total int

	query := `
		SELECT id, title, description, location, guests, beds, baths, price, cleaning_fee, service_fee, taxes, landlord_id, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err := r.db.SelectContext(r.ctx, &listings, query, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, 0, fmt.Errorf("error finding listings: %w", err)
	}

	totalQuery := "SELECT COUNT(*) FROM listings"
	err = r.db.GetContext(r.ctx, &total, totalQuery)

	if err != nil {
		return nil, 0, 0, err
	}

	totalPage := (total + limit - 1) / limit

	return listings, total, totalPage, nil
}

func (r *listingRepository) FindOne(id int) (*domain.Listing, error) {
	query := `
		SELECT id, title, description, location, guests, beds, baths, price, cleaning_fee, service_fee, taxes, landlord_id, created_at, updated_at
		FROM listings
		WHERE id = $1
	`

	listing := &domain.Listing{}

	err := r.db.GetContext(r.ctx, listing, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing with id %d not found", id)
		}
		return nil, fmt.Errorf("error finding listing: %w", err)
	}

	return listing, nil
}

func (r *listingRepository) Update(id int, listing *domain.Listing) (*domain.Listing, error) {
	query := `
		UPDATE listings
		SET title = $1, description = $2, location = $3, guests = $4, beds = $5, baths = $6, price = $7, cleaning_fee = $8, service_fee = $9, taxes = $10, landlord_id = $11, updated_at = $12
		WHERE id = $13
		RETURNING id, updated_at
		`
	now := time.Now()
	listing.UpdatedAt = now

	err := r.db.QueryRowxContext(r.ctx, query,
		listing.Title,
		listing.Description,
		listing.Location,
		listing.Guests,
		listing.Beds,
		listing.Baths,
		listing.Price,
		listing.CleaningFee,
		listing.ServiceFee,
		listing.Taxes,
		listing.LandlordID,
		listing.UpdatedAt,
		listing.ID,
	).Scan(&listing.ID, &listing.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating listing: %w", err)
	}

	return listing, err
}

func (r *listingRepository) Remove(id int) error {
	query := "DELETE FROM listings WHERE id = $1"

	_, err := r.db.ExecContext(r.ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting listing: %w", err)
	}

	return nil
}

func (r *listingRepository) Save(listing *domain.Listing) (*domain.Listing, error) {
	query := `
		INSERT INTO listings (title, description, location, guests, beds, baths, price, cleaning_fee, service_fee, taxes, landlord_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	listing.CreatedAt = now
	listing.UpdatedAt = now

	err := r.db.QueryRowxContext(r.ctx, query,
		listing.Title,
		listing.Description,
		listing.Location,
		listing.Guests,
		listing.Beds,
		listing.Baths,
		listing.Price,
		listing.CleaningFee,
		listing.ServiceFee,
		listing.Taxes,
		listing.LandlordID,
		listing.CreatedAt,
		listing.UpdatedAt,
	).Scan(&listing.ID, &listing.CreatedAt, &listing.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting listing: %w", err)
	}

	return listing, nil
}
