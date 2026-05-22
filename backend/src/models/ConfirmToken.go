package models

import "time"

type ConfirmToken struct {
    TokenID   string    `db:"token_id"`
    UserID    string    `db:"user_id"`
    Token     string    `db:"token"`
    Action    string    `db:"action"`
    ExpiresAt time.Time `db:"expires_at"`
    CreatedAt time.Time `db:"created_at"`
}