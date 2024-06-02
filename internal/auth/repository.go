package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// SessionRepository interface defines the methods that the repository layer will implement.
type dbRepository struct {
	db *sqlx.DB
}

// NewSessionRepository creates a new SessionRepository with a given database connection.
func NewSessionRepository(dbConnection *sqlx.DB) SessionRepository {
	return &dbRepository{
		db: dbConnection,
	}
}

// Gets a single session in the database.
func (r *dbRepository) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	session := &Session{}

	err := r.db.GetContext(ctx, session, "SELECT * FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Gets a single session in the database by refresh token.
func (r *dbRepository) GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	session := &Session{}

	err := r.db.GetContext(ctx, session, "SELECT * FROM sessions WHERE refresh_token = ?", refreshToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Gets all sessions for a user in the database.
func (r *dbRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) (*[]Session, error) {
	sessions := &[]Session{}

	// TODO: We need to take into account that a malicious user could try to create a lot of sessions for a single user.
	// making the query very slow.
	err := r.db.SelectContext(ctx, sessions, "SELECT * FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// Creates a new session in the database.
func (r *dbRepository) CreateSession(ctx context.Context, session *Session) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO sessions (id, refresh_token, user_id, type, ip, user_agent, created_at, last_used_at, expires_at) VALUES (:id, :refresh_token, :user_id, :type, :ip, :user_agent, :created_at, :last_used_at, :expires_at)", session)
	if err != nil {
		return err
	}

	return nil
}

// Updates a session in the database.
func (r *dbRepository) UpdateSession(ctx context.Context, session *Session) error {
	_, err := r.db.NamedExecContext(ctx, "UPDATE sessions SET refresh_token = :refresh_token, user_id = :user_id, type = :type, ip = :ip, user_agent = :user_agent, created_at = :created_at, last_used_at = :last_used_at, expires_at = :expires_at WHERE id = :id", session)
	if err != nil {
		return err
	}

	return nil
}

// Deletes a session in the database.
func (r *dbRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET expires_at = ? WHERE id = ?", now, sessionID)
	if err != nil {
		return err
	}

	return nil
}

// Deletes all sessions for a user in the database.
func (r *dbRepository) DeleteAllSessions(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET expires_at = ? WHERE user_id = ?", now, userID)
	if err != nil {
		return err
	}

	return nil
}
