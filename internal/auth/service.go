package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/internal/utils/random"
)

// Session service implementation.
type sessionService struct {
	sessionRepository SessionRepository
}

// Create a new 'service' or 'use-case' for 'Session' entity.
func NewSessionService(r SessionRepository) SessionService {
	return &sessionService{
		sessionRepository: r,
	}
}

// Get a single session by its ID.
func (s *sessionService) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	return s.sessionRepository.GetSession(ctx, sessionID)
}

// Get a single session by its refresh token.
func (s *sessionService) GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	return s.sessionRepository.GetSessionByToken(ctx, refreshToken)
}

// Get all sessions for a user.
func (s *sessionService) GetUserSessions(ctx context.Context, userID uuid.UUID) (*[]Session, error) {
	return s.sessionRepository.GetUserSessions(ctx, userID)
}

// Create a new session.
func (s *sessionService) CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent string, tokenType TokenType) (*Session, error) {
	rToken, err := random.GenerateRandomString(144)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		ID:           uuid.New(),
		UserID:       userID,
		IP:           ip,
		UserAgent:    userAgent,
		Type:         tokenType,
		RefreshToken: rToken,
		CreatedAt:    now,
		LastUsedAt:   now,
		ExpiresAt:    now.Add(90 * 24 * time.Hour),
	}

	err = s.sessionRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Refresh a session by its refresh token.
func (s *sessionService) RefreshSession(ctx context.Context, refreshToken, ip string) (*Session, error) {
	session, err := s.sessionRepository.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, consts.ErrSessionExpired
	}

	session.RefreshToken, err = random.GenerateRandomString(144)
	if err != nil {
		return nil, err
	}
	session.IP = ip
	session.LastUsedAt = time.Now()

	err = s.sessionRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Revoke a session by its ID.
func (s *sessionService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepository.DeleteSession(ctx, sessionID)
}

// Revoke a session by its refresh token.
func (s *sessionService) RevokeSessionByToken(ctx context.Context, refreshToken string) error {
	session, err := s.sessionRepository.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return s.sessionRepository.DeleteSession(ctx, session.ID)
}

// Revoke all sessions for a user.
func (s *sessionService) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	return s.sessionRepository.DeleteAllSessions(ctx, userID)
}
