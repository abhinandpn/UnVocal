package model

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Number   string `json:"number"`
	UserCode string `json:"user_code"`
}
