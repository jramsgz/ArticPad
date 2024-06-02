package user

import (
	"context"
	"database/sql"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/pkg/argon2id"
	"github.com/jramsgz/articpad/pkg/validator"
)

// User service implementation.
type userService struct {
	userRepository UserRepository
}

// Create a new 'service' or 'use-case' for 'User' entity.
func NewUserService(r UserRepository) UserService {
	return &userService{
		userRepository: r,
	}
}

// Get a user by ID.
func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.userRepository.GetUser(ctx, userID)
}

// Get a user by email.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}

// Get a user by username.
func (s *userService) GetUserByUsername(ctx context.Context, userName string) (*User, error) {
	return s.userRepository.GetUserByUsername(ctx, userName)
}

// Create a new user in the system. The user is validated before being created.
// The first user created in the system is an admin.
// ID, CreatedAt, UpdatedAt and DeletedAt are set automatically, providen values are ignored.
func (s *userService) CreateUser(ctx context.Context, user *User) error {
	err := s.validateUser(user)
	if err != nil {
		return err
	}

	foundUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && err != consts.ErrRecordNotFound {
		if err == consts.ErrDeletedRecord {
			return consts.ErrEmailDeactivated
		}
		return err
	}
	if foundUser != nil {
		return consts.ErrEmailAlreadyExists
	}

	foundUser, err = s.GetUserByUsername(ctx, user.Username)
	if err != nil && err != consts.ErrRecordNotFound {
		if err == consts.ErrDeletedRecord {
			return consts.ErrUsernameDeactivated
		}
		return err
	}
	if foundUser != nil {
		return consts.ErrUsernameAlreadyExists
	}

	if ok, err := s.IsFirstUser(ctx); ok && err == nil {
		user.IsAdmin = true
	}

	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	user.ID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	return s.userRepository.CreateUser(ctx, user)
}

// Update a user in the system. The user is validated before being updated.
// The email and username must be unique. The password is hashed before being updated if it has changed.
// ID, CreatedAt, UpdatedAt and DeletedAt are automatically taken care of so provided values are ignored.
func (s *userService) UpdateUser(ctx context.Context, user *User) error {
	err := s.validateUser(user)
	if err != nil {
		return err
	}

	foundUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && err != consts.ErrRecordNotFound {
		if err == consts.ErrDeletedRecord {
			return consts.ErrEmailDeactivated
		}
		return err
	}
	if foundUser != nil {
		return consts.ErrEmailAlreadyExists
	}

	foundUser, err = s.GetUserByUsername(ctx, user.Username)
	if err != nil && err != consts.ErrRecordNotFound {
		if err == consts.ErrDeletedRecord {
			return consts.ErrUsernameDeactivated
		}
		return err
	}
	if foundUser != nil {
		return consts.ErrUsernameAlreadyExists
	}

	actualUser, err := s.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	if actualUser.Password != user.Password {
		hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	user.UpdatedAt = time.Now()
	return s.userRepository.UpdateUser(ctx, user)
}

// Delete a user from the system.
func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepository.DeleteUser(ctx, userID)
}

// Check if the user is the first user in the system.
func (s *userService) IsFirstUser(ctx context.Context) (bool, error) {
	_, err := s.userRepository.GetFirstUser(ctx)
	if err != nil && err == consts.ErrRecordNotFound {
		return true, nil
	}

	return false, err
}

// Get a user by email or username. The user is searched first by username and then by email.
func (s *userService) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, emailOrUsername)
	if err != nil && err != consts.ErrRecordNotFound {
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

// Verify the user email. The user is verified if the verification token is valid and the user has not been verified before.
func (s *userService) VerifyUser(ctx context.Context, verificationToken string) error {
	user, err := s.userRepository.GetUserByVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}

	if user.VerifiedAt.Valid && user.VerifiedAt.Time.Before(time.Now()) {
		return consts.ErrEmailAlreadyVerified
	}

	user.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err = s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// Generate a password reset token for the user. The token is returned in the user object.
func (s *userService) GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.PasswordResetToken = sql.NullString{
		String: uuid.New().String(),
		Valid:  true,
	}

	user.PasswordResetExpiresAt = sql.NullTime{
		Time:  time.Now().Add(time.Hour * 4),
		Valid: true,
	}

	return user, s.userRepository.UpdateUser(ctx, user)
}

// Reset the user password. A password reset token is required to reset the password.
func (s *userService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	user, err := s.userRepository.GetUserByPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	if !user.PasswordResetExpiresAt.Valid || user.PasswordResetExpiresAt.Time.Before(time.Now()) {
		return consts.ErrPasswordResetTokenExpired
	}

	user.Password = newPassword
	user.PasswordResetExpiresAt.Time = time.Now()

	return s.userRepository.UpdateUser(ctx, user)
}

// Validates the user data and returns an error if it is not valid.
func (s *userService) validateUser(user *User) error {
	parsedEmail, err := mail.ParseAddress(user.Email)
	if err != nil || len(parsedEmail.Address) > 100 {
		return consts.ErrInvalidEmail
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
