package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represents the 'User' object.
type User struct {
	ID        uuid.UUID      `json:"ID" gorm:"primaryKey;type:uuid"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"password"`
	Verified  bool           `gorm:"not null" json:"verified"`
	Admin     bool           `gorm:"not null" json:"admin"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// BeforeCreate will set default values for the user.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.New()
	// Set the created and updated times.
	now := time.Now()
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
	GetUser(ctx context.Context, userID int) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, userID int, user *User) error
	DeleteUser(ctx context.Context, userID int) error
	GetFirstUser(ctx context.Context) (*User, error)
}

// Our use-case or service will implement these methods.
type UserService interface {
	GetUsers(ctx context.Context) (*[]User, error)
	GetUser(ctx context.Context, userID int) (*User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, userID int, user *User) error
	DeleteUser(ctx context.Context, userID int) error
	IsFirstUser(ctx context.Context) (bool, error)
}
