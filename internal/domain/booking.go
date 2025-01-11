package domain

import "time"

type Booking struct {
	ID            int       `json:"id" db:"id"`
	ListingID     int       `json:"-" db:"listing_id"`
	GuestID       int       `json:"guest_id" db:"guest_id"`
	Guests        int       `json:"guests" db:"guests"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	Nights        int       `json:"nights" db:"nights"`
	PhoneNumber   *string   `json:"phone_number" db:"phone_number"`
	MessageToHost *string   `json:"message_to_host" db:"message_to_host"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Payment struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	BookingID    int       `json:"-" db:"booking_id"`
	IsSuccessful bool      `json:"is_successful" db:"is_successful"`
	Price        float64   `json:"price" db:"price"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type PriceDetail struct {
	ID             int       `json:"id" db:"id"`
	BookingID      int       `json:"-" db:"booking_id"`
	TotalHomePrice float64   `json:"total_home_price" db:"total_home_price"`
	CleaningFee    float64   `json:"cleaning_fee" db:"cleaning_fee"`
	ServiceFee     float64   `json:"service_fee" db:"service_fee"`
	Taxes          float64   `json:"taxes" db:"taxes"`
	TotalPrice     float64   `json:"total_price" db:"total_price"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
