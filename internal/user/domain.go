package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represents the 'User' object.
type User struct {
	ID                     uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid"`
	Username               string         `json:"username" gorm:"uniqueIndex;not null"`
	Email                  string         `json:"email" gorm:"uniqueIndex;not null"`
	Password               string         `json:"-" gorm:"not null"`
	VerifiedAt             sql.NullTime   `json:"verified_at,omitempty"`
	VerificationToken      string         `json:"-" gorm:"uniqueIndex;not null"`
	PasswordResetToken     sql.NullString `json:"-" gorm:"uniqueIndex"`
	PasswordResetExpiresAt time.Time      `json:"-"`
	IsAdmin                bool           `json:"is_admin" gorm:"not null"`
	Lang                   string         `json:"lang" gorm:"not null"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate will set default values for the user.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.New()
	// Set the created and updated times.
	now := time.Now()
	user.PasswordResetExpiresAt = now
	user.CreatedAt = now
	user.UpdatedAt = now
	return
}

// BeforeUpdate will set default values for the user.
func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}

// Our repository will implement these methods.
type UserRepository interface {
	GetUsers(ctx context.Context) (*[]User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	GetFirstUser(ctx context.Context) (*User, error)
	VerifyUser(ctx context.Context, verificationToken string) error
	SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetUserByPasswordResetToken(ctx context.Context, token string) (*User, error)
}

// Our use-case or service will implement these methods.
type UserService interface {
	GetUsers(ctx context.Context) (*[]User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	IsFirstUser(ctx context.Context) (bool, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error)
	VerifyUser(ctx context.Context, verificationToken string) error
	SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	ResetPassword(ctx context.Context, token string, password string) error
}
