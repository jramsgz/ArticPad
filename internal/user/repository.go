package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jramsgz/articpad/internal/utils/consts"
)

// UserRepository interface defines the methods that the repository layer will implement.
type dbRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository with a given database connection.
func NewUserRepository(dbConnection *sqlx.DB) UserRepository {
	return &dbRepository{
		db: dbConnection,
	}
}

// Gets a single user in the database.
func (r *dbRepository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, consts.ErrDeletedRecord
	}

	return user, nil
}

// Gets a single user in the database by email.
func (r *dbRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, consts.ErrDeletedRecord
	}

	return user, nil
}

// Gets a single user in the database by username.
func (r *dbRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, consts.ErrDeletedRecord
	}

	return user, nil
}

// Creates a single user in the database.
func (r *dbRepository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, username, email, password, verification_token, is_admin, lang, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Email, user.Password, user.VerificationToken, user.IsAdmin, user.Lang, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a single user in the database.
func (r *dbRepository) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET username = ?, email = ?, password = ?, verified_at = ?, verification_token = ?, password_reset_token = ?, password_reset_expires_at = ?, is_admin = ?, lang = ?, updated_at = ? WHERE id = ?",
		user.Username, user.Email, user.Password, user.VerifiedAt, user.VerificationToken, user.PasswordResetToken, user.PasswordResetExpiresAt, user.IsAdmin, user.Lang, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// Deletes a single user in the database.
func (r *dbRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, "UPDATE users SET deleted_at = ? WHERE id = ?", now, userID)
	if err != nil {
		return err
	}

	return nil
}

// Gets the first user in the database.
func (r *dbRepository) GetFirstUser(ctx context.Context) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users ORDER BY created_at ASC LIMIT 1")
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByVerificationToken returns a user by its verification token.
func (r *dbRepository) GetUserByVerificationToken(ctx context.Context, verificationToken string) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE verification_token = ?", verificationToken)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, consts.ErrDeletedRecord
	}

	return user, nil
}

// Gets a user by its password reset token and checks if it is still valid.
func (r *dbRepository) GetUserByPasswordResetToken(ctx context.Context, token string) (*User, error) {
	user := &User{}

	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE password_reset_token = ?", token)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt.Valid {
		return nil, consts.ErrDeletedRecord
	}

	return user, nil
}
