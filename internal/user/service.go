package user

import (
	"context"
	"time"
)

// Implementation of the repository in this service.
type userService struct {
	userRepository UserRepository
}

// Create a new 'service' or 'use-case' for 'User' entity.
func NewUserService(r UserRepository) UserService {
	return &userService{
		userRepository: r,
	}
}

// Implementation of 'GetUsers'.
func (s *userService) GetUsers(ctx context.Context) (*[]User, error) {
	return s.userRepository.GetUsers(ctx)
}

// Implementation of 'GetUser'.
func (s *userService) GetUser(ctx context.Context, userID int) (*User, error) {
	return s.userRepository.GetUser(ctx, userID)
}

// Implementation of 'GetUserByEmail'.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}

// Implementation of 'GetUserByUsername'.
func (s *userService) GetUserByUsername(ctx context.Context, userName string) (*User, error) {
	return s.userRepository.GetUserByUsername(ctx, userName)
}

// Implementation of 'CreateUser'.
func (s *userService) CreateUser(ctx context.Context, user *User) error {
	// Set default value of 'CreatedAt' and 'UpdatedAt'.
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Pass to the repository layer.
	return s.userRepository.CreateUser(ctx, user)
}

// Implementation of 'UpdateUser'.
func (s *userService) UpdateUser(ctx context.Context, userID int, user *User) error {
	// Set value for 'UpdatedAt' attribute.
	user.UpdatedAt = time.Now()

	// Pass to the repository layer.
	return s.userRepository.UpdateUser(ctx, userID, user)
}

// Implementation of 'DeleteUser'.
func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	return s.userRepository.DeleteUser(ctx, userID)
}

// Implementation of 'IsFirstUser'.
func (s *userService) IsFirstUser(ctx context.Context) (bool, error) {
	_, err := s.userRepository.GetFirstUser(ctx)
	if err != nil && err.Error() == "record not found" {
		return true, nil
	}

	return false, err
}
