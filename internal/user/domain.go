package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jramsgz/articpad/pkg/argon2id"
)

// Represents the 'User' object.
type User struct {
	ID                     uuid.UUID      `db:"id"                        json:"id"`
	Username               string         `db:"username"                  json:"username"`
	Email                  string         `db:"email"                     json:"email"`
	Password               string         `db:"password"                  json:"-"`
	VerifiedAt             sql.NullTime   `db:"verified_at"               json:"-"`
	VerificationToken      string         `db:"verification_token"        json:"-"`
	PasswordResetToken     sql.NullString `db:"password_reset_token"      json:"-"`
	PasswordResetExpiresAt sql.NullTime   `db:"password_reset_expires_at" json:"-"`
	IsAdmin                bool           `db:"is_admin"                  json:"is_admin"`
	Lang                   string         `db:"lang"                      json:"lang"`
	CreatedAt              time.Time      `db:"created_at"                json:"created_at"`
	UpdatedAt              time.Time      `db:"updated_at"                json:"updated_at"`
	DeletedAt              sql.NullTime   `db:"deleted_at"                json:"-"`
}

// Checks if provided password matches the user's password.
func (u *User) ComparePassword(password string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, u.Password)
}

// UserRepository is an interface that defines the methods that the user repository must implement.
type UserRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByVerificationToken(ctx context.Context, verificationToken string) (*User, error)
	GetUserByPasswordResetToken(ctx context.Context, token string) (*User, error)
	GetFirstUser(ctx context.Context) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

// UserService is an interface that defines the methods that the user service must implement.
type UserService interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	IsFirstUser(ctx context.Context) (bool, error)
	VerifyUser(ctx context.Context, verificationToken string) error
	GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (*User, error)
	ResetPassword(ctx context.Context, token string, password string) error
}
