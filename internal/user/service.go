package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
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
func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error {
	// Set value for 'UpdatedAt' attribute.
	user.UpdatedAt = time.Now()

	// Pass to the repository layer.
	return s.userRepository.UpdateUser(ctx, userID, user)
}

// Implementation of 'DeleteUser'.
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepository.DeleteUser(ctx, userID)
}

// Implementation of 'IsFirstUser'.
func (s *userService) IsFirstUser(ctx context.Context) (bool, error) {
	_, err := s.userRepository.GetFirstUser(ctx)
	if err != nil && err == gorm.ErrRecordNotFound {
		return true, nil
	}

	return false, err
}

// Implementation of 'GetUserByEmailOrUsername'.
func (s *userService) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, emailOrUsername)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == nil {
		return user, nil
	}

	user, err = s.userRepository.GetUserByEmail(ctx, emailOrUsername)
	if err != nil {
		return nil, err
	}

	return user, err
}

// Implementation of 'VerifyUser'.
func (s *userService) VerifyUser(ctx context.Context, verificationToken string) error {
	return s.userRepository.VerifyUser(ctx, verificationToken)
}

// Implementation of 'SetPasswordResetToken'.
func (s *userService) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	return s.userRepository.SetPasswordResetToken(ctx, userID, token, expiresAt)
}

// Implementation of 'ResetPassword'.
func (s *userService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	user, err := s.userRepository.GetUserByPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Set new password.
	user.Password = newPassword
	user.UpdatedAt = time.Now()

	// Pass to the repository layer.
	return s.userRepository.UpdateUser(ctx, user.ID, user)
}
