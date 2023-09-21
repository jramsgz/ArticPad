package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
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
	// Initialize variables.
	var users []User

	// Get all users.
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return all of our users.
	return &users, nil
}

// Gets a single user in the database.
func (r *dbRepository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).First(user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return result.
	return user, nil
}

// Gets a single user in the database by email.
func (r *dbRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return result.
	return user, nil
}

// Gets a single user in the database by username.
func (r *dbRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return result.
	return user, nil
}

// Creates a single user in the database.
func (r *dbRepository) CreateUser(ctx context.Context, user *User) error {
	// Insert one user.
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Updates a single user in the database.
func (r *dbRepository) UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error {
	// Update one user.
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Deletes a single user in the database.
func (r *dbRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Delete one user.
	result := r.db.WithContext(ctx).Delete(&User{}, userID)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Gets the first user in the database.
func (r *dbRepository) GetFirstUser(ctx context.Context) (*User, error) {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return result.
	return user, nil
}

// Verifies a user given its verification token.
func (r *dbRepository) VerifyUser(ctx context.Context, verificationToken string) error {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).Where("verification_token = ?", verificationToken).First(user)
	if result.Error != nil {
		return result.Error
	}

	// Check if user is already verified.
	if user.VerifiedAt.Valid && user.VerifiedAt.Time.Before(time.Now()) {
		return errors.New("user is already verified")
	}

	// Update user.
	user.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	// Save user.
	result = r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Sets the password reset token for a user.
func (r *dbRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).Where("id = ?", userID).First(user)
	if result.Error != nil {
		return result.Error
	}

	// Update user.
	user.PasswordResetToken = sql.NullString{
		String: token,
		Valid:  true,
	}
	user.PasswordResetExpiresAt = expiresAt

	// Save user.
	result = r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Gets a user by its password reset token and checks if it is still valid.
func (r *dbRepository) GetUserByPasswordResetToken(ctx context.Context, token string) (*User, error) {
	// Initialize variable.
	user := &User{}

	// Prepare SQL to get one user.
	result := r.db.WithContext(ctx).Where("password_reset_token = ?", token).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if token is still valid.
	if user.PasswordResetExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// Return result.
	return user, nil
}
