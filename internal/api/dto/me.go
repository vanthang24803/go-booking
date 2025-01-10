package dto

type UpdateProfileRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,alpha"`
	Surname   string `json:"surname" validate:"required,alpha"`
}
