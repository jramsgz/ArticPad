package user

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/pkg/argon2id"
	"github.com/jramsgz/articpad/pkg/validator"
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
	err := s.validateUser(ctx, user)
	if err != nil {
		return err
	}

	foundUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err.Error() == consts.ErrDeletedRecord {
			return errors.New(consts.ErrEmailDeactivated)
		}
		return err
	}
	if foundUser != nil {
		return errors.New(consts.ErrEmailAlreadyExists)
	}

	foundUser, err = s.GetUserByUsername(ctx, user.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err.Error() == consts.ErrDeletedRecord {
			return errors.New(consts.ErrUsernameDeactivated)
		}
		return err
	}
	if foundUser != nil {
		return errors.New(consts.ErrUsernameAlreadyExists)
	}

	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepository.CreateUser(ctx, user)
}

// Implementation of 'UpdateUser'.
func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error {
	err := s.validateUser(ctx, user)
	if err != nil {
		return err
	}

	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

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
	user, err := s.userRepository.GetUserByVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}

	if user.VerifiedAt.Valid && user.VerifiedAt.Time.Before(time.Now()) {
		return errors.New(consts.ErrEmailAlreadyVerified)
	}

	err = s.userRepository.SetUserVerified(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
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

	if user.PasswordResetExpiresAt.Before(time.Now()) {
		return errors.New(consts.ErrPasswordResetTokenExpired)
	}

	user.Password = newPassword
	user.PasswordResetExpiresAt = time.Now()

	err = s.validateUser(ctx, user)
	if err != nil {
		return err
	}

	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepository.UpdateUser(ctx, user.ID, user)
}

// Validates the user data and returns an error if it is not valid.
func (s *userService) validateUser(ctx context.Context, user *User) error {
	parsedEmail, err := mail.ParseAddress(user.Email)
	if err != nil || (err == nil && len(parsedEmail.Address) > 100) {
		return errors.New(consts.ErrInvalidEmail)
	}
	user.Email = parsedEmail.Address

	usernameValidator := validator.DefaultUsernameValidator
	if err := usernameValidator.Validate(user.Username); err != nil {
		return err
	}

	passwordValidator := validator.DefaultPasswordValidator([]string{user.Username, user.Email})
	if err := passwordValidator.Validate(user.Password); err != nil {
		return err
	}

	return nil
}
