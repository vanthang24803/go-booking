package dto

type CatalogRequest struct {
	Name string `json:"name" validate:"required"`
}
