package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/codebasky/lenslocked/rand"
)

const (
	DefaultResetDuration = 1 * time.Hour
	ResetTokenLength     = 32
	TokenExpireDuration  = 2 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	db *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this value is not set or is less than the
	// MinBytesPerToken const it will be ignored and MinBytesPerToken will be
	// used.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (service *PasswordResetService) Create(email string, uid int) (*PasswordReset, error) {
	token, err := rand.String(service.BytesPerToken)
	if err != nil {
		return nil, err
	}
	prst := PasswordReset{
		UserID:    uid,
		Token:     token,
		TokenHash: rand.Hash(token),
		ExpiresAt: time.Now().Add(service.Duration),
	}
	row := service.db.QueryRow(`INSERT INTO password_resets (user_id, token_hash, expires_at)
								VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
								UPDATE SET token_hash = $2, expires_at = $3 RETURNING id;`,
		prst.UserID, prst.TokenHash, prst.ExpiresAt)
	err = row.Scan(&prst.ID)
	if err != nil {
		return nil, fmt.Errorf("insert password token failed: %v", err)
	}
	return &prst, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	var user User
	row := service.db.QueryRow(`SELECT u.email, u.password_hash, u.id
		FROM password_resets JOIN users as u ON password_resets.user_id = u.id
		WHERE password_resets.token_hash = $1;`,
		rand.Hash(token))
	err := row.Scan(&user.Email, &user.Password_Hash, &user.ID)
	if err != nil {
		return nil, fmt.Errorf("not able to find user for token: %v", err)
	}
	return &user, nil
}

func (service *PasswordResetService) Delete(token string) error {
	tokenHash := rand.Hash(token)
	_, err := service.db.Exec(`DELETE FROM password_resets WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func NewPwdResetService(db *sql.DB) *PasswordResetService {
	return &PasswordResetService{
		db:            db,
		BytesPerToken: ResetTokenLength,
		Duration:      TokenExpireDuration,
	}
}
