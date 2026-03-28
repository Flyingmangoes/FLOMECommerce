package models

import "time"

type RefreshToken struct {
	TokenID   string    `db:"token_id"`
    UserID    string    `db:"user_id"`
	Token     string    `db:"token"`		 
	ExpiresAt time.Time `db:"expires_at"`
    CreatedAt time.Time `db:"created_at"`
	Revoked   bool		`db:"is_revoked"`
}	