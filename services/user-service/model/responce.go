package model

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Number   string `json:"number"`
	UserCode string `json:"user_code"`
}

type LoginResponse struct {
	RefreshToken string        `json:"refresh_token"`
	AccessToken  string        `json:"access_token"`
	User         *UserResponse `json:"user"`
}
