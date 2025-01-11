package domain

import "time"

type Listing struct {
	ID          int      `json:"id" db:"id"`
	LandlordID  int      `json:"landlord_id" db:"landlord_id"`
	Title       string   `json:"title" db:"title"`
	Description string   `json:"description" db:"description"`
	Location    string   `json:"location" db:"location"`
	Guests      int      `json:"guests" db:"guests"`
	Beds        int      `json:"beds" db:"beds"`
	Baths       int      `json:"baths" db:"baths"`
	Price       float64  `json:"price" db:"price"`
	CleaningFee float64  `json:"cleaning_fee" db:"cleaning_fee"`
	ServiceFee  float64  `json:"service_fee" db:"service_fee"`
	Taxes       float64  `json:"taxes" db:"taxes"`
	Photos      []Photo  `json:"photos" db:"-"`
	Review      []Review `json:"reviews" db:"-"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Photo struct {
	ID        int       `json:"id" db:"id"`
	ListingID int       `json:"-" db:"listing_id"`
	PublicID  string    `json:"-" db:"public_id"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Review struct {
	ID          int       `json:"id" db:"id"`
	ListingID   int       `json:"-" db:"listing_id"`
	AuthorID    int       `json:"-" db:"author_id"`
	BookingID   int       `json:"-" db:"booking_id"`
	Rating      int       `json:"rating" db:"rating"`
	Comment     string    `json:"comment" db:"comment"`
	IsPublished bool      `json:"is_published" db:"is_published"`
	IsEdited    bool      `json:"is_edited" db:"is_edited"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
