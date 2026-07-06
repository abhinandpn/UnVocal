package model

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Number   string `json:"number" example:"9876543210"`
	Password string `json:"password" example:"StrongPassword123"`
}

type UpdateRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Number   string `json:"number" example:"9876543210"`
	Password string `json:"password" example:"StrongPassword123"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" example:"john@example.com"`
	Password   string `json:"password" example:"StrongPassword123"`
}
