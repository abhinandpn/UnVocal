package model

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Number    string `json:"number"`
	CreatedAt string `json:"created_at"`
	Password  string `json:"-"`
}
