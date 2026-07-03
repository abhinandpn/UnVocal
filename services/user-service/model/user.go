package model

import "time"

type User struct {
	ID        string     `json:"id"`
	UserCode  string     `json:"user_code"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Number    string     `json:"number"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
