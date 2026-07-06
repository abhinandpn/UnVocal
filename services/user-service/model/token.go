package model

import "time"

type RefreshToken struct {
	ID        int64      `gorm:"primaryKey"`
	UserCode  string     `gorm:"column:user_code"`
	Token     string     `gorm:"unique"`
	ExpiresAt time.Time  `gorm:"column:expires_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
}
