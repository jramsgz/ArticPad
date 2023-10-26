package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"gorm.io/gorm"
)

// Represents that we will use MariaDB in order to implement the methods.
type dbRepository struct {
	db *gorm.DB
}

// Create a new repository with MariaDB as the driver.
func NewUserRepository(dbConnection *gorm.DB) UserRepository {
	return &dbRepository{
		db: dbConnection,
	}
}

// Gets all users in the database.
func (r *dbRepository) GetUsers(ctx context.Context) (*[]User, error) {
	var users []User

	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return &users, nil
}

// Gets a single user in the database.
func (r *dbRepository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}

	result := r.db.WithContext(ctx).First(user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Gets a single user in the database by email.
func (r *dbRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	result := r.db.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Gets a single user in the database by username.
func (r *dbRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}

	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Creates a single user in the database.
func (r *dbRepository) CreateUser(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Updates a single user in the database.
func (r *dbRepository) UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error {
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Deletes a single user in the database.
func (r *dbRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&User{}, userID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Gets the first user in the database.
func (r *dbRepository) GetFirstUser(ctx context.Context) (*User, error) {
	user := &User{}

	result := r.db.WithContext(ctx).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Verifies a user given its verification token.
func (r *dbRepository) VerifyUser(ctx context.Context, verificationToken string) error {
	user := &User{}

	result := r.db.WithContext(ctx).Where("verification_token = ?", verificationToken).First(user)
	if result.Error != nil {
		return result.Error
	}

	if user.VerifiedAt.Valid && user.VerifiedAt.Time.Before(time.Now()) {
		return errors.New(consts.ErrEmailAlreadyVerified)
	}

	user.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	result = r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Sets the password reset token for a user.
func (r *dbRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	user := &User{}

	result := r.db.WithContext(ctx).Where("id = ?", userID).First(user)
	if result.Error != nil {
		return result.Error
	}

	user.PasswordResetToken = sql.NullString{
		String: token,
		Valid:  true,
	}
	user.PasswordResetExpiresAt = expiresAt

	result = r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Gets a user by its password reset token and checks if it is still valid.
func (r *dbRepository) GetUserByPasswordResetToken(ctx context.Context, token string) (*User, error) {
	user := &User{}

	result := r.db.WithContext(ctx).Where("password_reset_token = ?", token).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	if user.PasswordResetExpiresAt.Before(time.Now()) {
		return nil, errors.New(consts.ErrPasswordResetTokenExpired)
	}

	return user, nil
}
