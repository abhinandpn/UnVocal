package model

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Number   string `json:"number"`
	Password string `json:"password"`
}
