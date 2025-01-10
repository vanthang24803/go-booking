package dto

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"

	User    = "user"
	Admin   = "admin"
	Manager = "manager"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required,alpha"`
	Surname   string `json:"surname" validate:"required,alpha"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
