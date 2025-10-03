package models

import "time"

// Token represents a session token
type Token struct {
	ID        string    `json:"id" db:"id"`
	Token     string    `json:"token" db:"token"`
	UserID    string    `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
}
