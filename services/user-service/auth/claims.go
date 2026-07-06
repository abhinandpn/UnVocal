package auth

import jwt "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserCode string `json:"user_code"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
