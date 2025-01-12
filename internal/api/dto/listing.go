package dto

type ListingRequest struct {
	Title       string  `json:"title" db:"title"`
	Description string  `json:"description" db:"description"`
	Catalogs    []int   `json:"catalogs" db:"-"`
	Location    string  `json:"location" db:"location"`
	Guests      int     `json:"guests" db:"guests"`
	Beds        int     `json:"beds" db:"beds"`
	Baths       int     `json:"baths" db:"baths"`
	Price       float64 `json:"price" db:"price"`
	CleaningFee float64 `json:"cleaning_fee" db:"cleaning_fee"`
	ServiceFee  float64 `json:"service_fee" db:"service_fee"`
	Taxes       float64 `json:"taxes" db:"taxes"`
}
