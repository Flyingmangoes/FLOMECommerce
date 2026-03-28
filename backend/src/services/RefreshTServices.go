package services

import (
	"backend/src/models"
	"context"
	"time"
	"database/sql"
	"fmt"
)

type TokenStoreInterface interface {
	SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
    GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
    DeleteRefreshToken(ctx context.Context, token string) error
    DeleteAllUserTokens(ctx context.Context, userID string) error
}

type TokenStore struct {
	db *sql.DB
}

func NewTokenStore(db *sql.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (ts *TokenStore) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) (error) {
	_, err := ts.db.ExecContext(ctx,
		`INSERT INTO mkt_refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)

	if err != nil {
		return err
	}
	
	return nil
}

func (ts *TokenStore) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}

	err := ts.db.QueryRowContext(ctx,
		`SELECT token_id, user_id, token, expires_at, created_at
		FROM mkt_refresh_tokens
		WHERE token = $1`,
		token,
	).Scan(&rt.TokenID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)

	if err != nil {
		return nil, err
	}

	return rt, nil
}

func (ts *TokenStore) DeleteRefreshToken(ctx context.Context, token string) error {
	result, err := ts.db.ExecContext(ctx,
		`DELETE FROM mkt_refresh_tokens WHERE token = $1`, token,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("Token not found")
	}

	return nil
}

func (ts *TokenStore) DeleteAllUserTokens(ctx context.Context, userID string) error {
	result, err := ts.db.ExecContext(ctx,
		`DELETE FROM mkt_refresh_tokens WHERE user_id = $1`, userID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("Token not found")
	}

	return nil
}