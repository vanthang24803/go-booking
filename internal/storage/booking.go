package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type BookingRepository interface {
	Save(req *domain.Booking) (*domain.Booking, error)
	FindDetail(id int) (*domain.Booking, error)
	FindAllForListing(listingId int, page int, limit int) ([]*domain.Booking, int, int, error)
	FindAllForUser(userId int, page int, limit int) ([]*domain.Booking, int, int, error)
	ExistBooking(listingId int, startDate time.Time, endDate time.Time) (bool, error)
}

type bookingRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewBookingRepository(db *sqlx.DB, ctx context.Context) *bookingRepository {
	return &bookingRepository{db: db, ctx: ctx}
}

func (r *bookingRepository) FindAllForListing(listingId, page, limit int) ([]*domain.Booking, int, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, listing_id, user_id, start_date, end_date, guest, nights, phone_number, created_at, updated_at
		FROM bookings
		WHERE listing_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookings []*domain.Booking
	err := r.db.SelectContext(r.ctx, &bookings, query, listingId, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching bookings for listing: %w", err)
	}

	totalQuery := "SELECT COUNT(*) FROM bookings WHERE listing_id = $1"
	var totalBookings int
	err = r.db.GetContext(r.ctx, &totalBookings, totalQuery, listingId)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching total bookings count: %w", err)
	}

	totalPages := (totalBookings + limit - 1) / limit

	return bookings, totalBookings, totalPages, nil
}

func (r *bookingRepository) FindAllForUser(userId, page, limit int) ([]*domain.Booking, int, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, listing_id, user_id, start_date, end_date, guest, nights, phone_number, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var bookings []*domain.Booking
	err := r.db.SelectContext(r.ctx, &bookings, query, userId, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching bookings for user: %w", err)
	}

	totalQuery := "SELECT COUNT(*) FROM bookings WHERE user_id = $1"
	var totalBookings int
	err = r.db.GetContext(r.ctx, &totalBookings, totalQuery, userId)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error fetching total bookings count: %w", err)
	}

	totalPages := (totalBookings + limit - 1) / limit

	return bookings, totalBookings, totalPages, nil
}

func (r *bookingRepository) FindDetail(id int) (*domain.Booking, error) {
	var booking *domain.Booking

	query := `SELECT id, listing_id, user_id, start_date, end_date, guest, nights, phone_number, message_to_host, created_at, updated_at FROM bookings WHERE id = $1`

	if err := r.db.GetContext(r.ctx, &booking, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking with id %d not found", id)
		}
		return nil, fmt.Errorf("error finding booking: %w", err)
	}

	return booking, nil
}

func (r *bookingRepository) Save(booking *domain.Booking) (*domain.Booking, error) {
	query := `
		INSERT INTO bookings (listing_id, user_id, start_date, end_date, guest, nights, phone_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	err := r.db.QueryRowxContext(
		r.ctx,
		query,
		booking.ListingID,
		booking.GuestID,
		booking.StartDate,
		booking.EndDate,
		booking.Guests,
		booking.Nights,
		booking.PhoneNumber,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error saving booking: %w", err)
	}

	return booking, nil
}

func (r *bookingRepository) ExistBooking(listingId int, startDate time.Time, endDate time.Time) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM bookings WHERE listing_id = $1 AND start_date <= $2 AND end_date >= $3)"
	var exists bool
	err := r.db.GetContext(r.ctx, &exists, query, listingId, startDate, endDate)
	if err != nil {
		return false, fmt.Errorf("error checking for existing booking: %w", err)
	}
	return exists, nil
}
