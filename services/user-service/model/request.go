package model

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Number   string `json:"number" example:"9876543210"`
	Password string `json:"password" example:"StrongPassword123"`
}
