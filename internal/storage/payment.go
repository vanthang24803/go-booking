package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/may20xx/booking/internal/domain"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) (*domain.Payment, error)
	FindForBooking(bookingId int) (*domain.Payment, error)
}

type paymentRepository struct {
	db  *sqlx.DB
	ctx context.Context
}

func NewPaymentRepository(db *sqlx.DB, ctx context.Context) *paymentRepository {
	return &paymentRepository{db: db, ctx: ctx}
}

func (r *paymentRepository) Save(payment *domain.Payment) (*domain.Payment, error) {

	query := `
		INSERT INTO payments (booking_id, name, is_successful, price, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()

	payment.CreatedAt = now

	err := r.db.QueryRowxContext(r.ctx, query,
		payment.BookingID,
		payment.Name,
		payment.IsSuccessful,
		payment.Price,
		payment.CreatedAt,
	).Scan(&payment.ID, &payment.CreatedAt)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) FindForBooking(bookingId int) (*domain.Payment, error) {
	var payment *domain.Payment

	query := `
		SELECT id, booking_id, name, is_successful, price, created_at
		FROM payments
		WHERE booking_id = $1
	`

	if err := r.db.GetContext(r.ctx, &payment, query, bookingId); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment for booking with booking_id id %d not found", bookingId)
		}
		return nil, fmt.Errorf("error finding payment for booking: %w", err)
	}

	return payment, nil
}
