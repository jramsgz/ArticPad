package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	User TokenType = "user"
	API  TokenType = "api"
)

// Represents the 'Session' object.
type Session struct {
	ID           uuid.UUID `db:"id"            json:"id"`
	RefreshToken string    `db:"refresh_token" json:"-"`
	UserID       uuid.UUID `db:"user_id"       json:"user_id"`
	Type         TokenType `db:"type"          json:"type"`
	IP           string    `db:"ip"            json:"ip"`
	UserAgent    string    `db:"user_agent"    json:"browser"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	LastUsedAt   time.Time `db:"last_used_at"  json:"last_used_at"`
	ExpiresAt    time.Time `db:"expires_at"    json:"expires_at"`
}

// SessionRepository interface defines the methods that the repository layer will implement.
type SessionRepository interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) (*[]Session, error)
	CreateSession(ctx context.Context, session *Session) error
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteAllSessions(ctx context.Context, userID uuid.UUID) error
}

// SessionService interface defines the methods that the service layer will implement.
type SessionService interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) (*[]Session, error)
	CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent string, tokenType TokenType) (*Session, error)
	RefreshSession(ctx context.Context, refreshToken, ip string) (*Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeSessionByToken(ctx context.Context, refreshToken string) error
	RevokeAllSessions(ctx context.Context, userID uuid.UUID) error
}
