package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"rest-api-golang/models"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// CreateToken creates a new token
func (r *TokenRepository) CreateToken(token *models.Token) error {
	query := `
		INSERT INTO tokens (id, token, user_id, expires_at, created_at, is_revoked)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		token.ID,
		token.Token,
		token.UserID,
		token.ExpiresAt,
		token.CreatedAt,
		token.IsRevoked,
	)

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}

// GetTokenByValue retrieves a token by its value
func (r *TokenRepository) GetTokenByValue(tokenValue string) (*models.Token, error) {
	query := `
		SELECT id, token, user_id, expires_at, created_at, is_revoked
		FROM tokens 
		WHERE token = $1 AND is_revoked = false AND expires_at > $2`

	token := &models.Token{}
	err := r.db.QueryRow(query, tokenValue, time.Now()).Scan(
		&token.ID,
		&token.Token,
		&token.UserID,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.IsRevoked,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}

// RevokeToken marks a token as revoked
func (r *TokenRepository) RevokeToken(tokenValue string) error {
	query := `UPDATE tokens SET is_revoked = true WHERE token = $1`

	result, err := r.db.Exec(query, tokenValue)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens
func (r *TokenRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM tokens WHERE expires_at < $1`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}
