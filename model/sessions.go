package model

import (
	"database/sql"
	"fmt"

	"github.com/codebasky/lenslocked/rand"
)

const (
	SessonTokenLength = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	db *sql.DB
}

func NewSessionSrv(db *sql.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

func (ss *SessionService) Create(uid int) (*Session, error) {
	token, err := rand.String(SessonTokenLength)
	if err != nil {
		return nil, err
	}
	tokenHash := rand.Hash(token)

	s := Session{
		UserID:    uid,
		Token:     token,
		TokenHash: tokenHash,
	}

	row := ss.db.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE SET token_hash = $2 RETURNING id;`,
		s.UserID, s.TokenHash)
	err = row.Scan(&s.ID)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}
	return &s, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := rand.Hash(token)
	var user User

	row := ss.db.QueryRow(`
		SELECT u.email, u.password_hash
		FROM sessions JOIN users as u ON sessions.user_id = u.id
		WHERE sessions.token_hash = $1;`,
		tokenHash)
	err := row.Scan(&user.Email, &user.Password_Hash)
	if err != nil {
		return nil, fmt.Errorf("not able to find user for token: %v", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := rand.Hash(token)
	_, err := ss.db.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
