package domain

type Catalog struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CatalogListing struct {
	CatalogID int `db:"catalog_id"`
	ListingID int `db:"listing_id"`
}
