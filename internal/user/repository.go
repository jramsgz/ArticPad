package user

import (
	"context"

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
func (r *dbRepository) GetUser(ctx context.Context, userID int) (*User, error) {
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
func (r *dbRepository) UpdateUser(ctx context.Context, userID int, user *User) error {
	// Update one user.
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}

// Deletes a single user in the database.
func (r *dbRepository) DeleteUser(ctx context.Context, userID int) error {
	// Delete one user.
	result := r.db.WithContext(ctx).Delete(&User{}, userID)
	if result.Error != nil {
		return result.Error
	}

	// Return empty.
	return nil
}
