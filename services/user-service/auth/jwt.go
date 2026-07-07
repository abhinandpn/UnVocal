package auth

import (
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 30 * 24 * time.Hour
	Issuer             = "unvocal-user-service"
)

func generateToken(userCode, secret string, expiry time.Duration) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(expiry)

	claims := Claims{
		UserCode: userCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // Unique JWT ID (jti)
			Subject:   userCode,
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateAccessToken creates a short-lived JWT.
func GenerateAccessToken(userCode string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}

	token, _, err := generateToken(userCode, secret, AccessTokenExpiry)
	return token, err
}

// GenerateRefreshToken creates a long-lived JWT.
func GenerateRefreshToken(userCode string) (string, time.Time, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return "", time.Time{}, errors.New("JWT_REFRESH_SECRET is not configured")
	}

	return generateToken(userCode, secret, RefreshTokenExpiry)
}

// ValidateAccessToken validates an access token.
func ValidateAccessToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, os.Getenv("JWT_SECRET"))
}

// ValidateRefreshToken validates a refresh token.
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, os.Getenv("JWT_REFRESH_SECRET"))
}

func validateToken(tokenString, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("JWT secret is not configured")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
