package auth

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/mailer"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/internal/utils/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Create an authentication handler.
type AuthHandler struct {
	userService user.UserService
	mailer      *mailer.Mailer
}

// Creates a new authentication handler.
func NewAuthHandler(authRoute fiber.Router, us user.UserService, mail *mailer.Mailer) {
	handler := &AuthHandler{
		userService: us,
		mailer:      mail,
	}

	// Declare routing for specific routes.
	authRoute.Post("/login", handler.signInUser)
	authRoute.Post("/register", handler.signUpUser)
	authRoute.Post("/logout", JWTMiddleware(), handler.logOutUser) // TODO
	authRoute.Post("/resend", handler.resendVerificationEmail)
	authRoute.Get("/verify/:token", handler.verifyUser)
	authRoute.Post("/forgot", handler.forgotPassword)
	authRoute.Get("/reset", handler.resetPassword)
	authRoute.Get("/refresh", JWTMiddleware(), handler.refreshToken) // TODO
	authRoute.Get("/me", JWTMiddleware(), handler.getMe)             // TODO
}

// checkPasswordHash compare password with hash
func (h *AuthHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Signs in a user and gives them a JWT.
func (h *AuthHandler) signInUser(c *fiber.Ctx) error {
	// Create a struct so the request body can be mapped here.
	type loginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	// Create a struct for our custom JWT payload.
	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		UserIP string `json:"user_ip"`
		jwt.RegisteredClaims
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get request body.
	request := &loginRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Get user by username from database.
	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidCredentials)
	} else if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// If password is incorrect, do not allow access.
	if !h.checkPasswordHash(request.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidCredentials)
	}

	if config.GetString("ENABLE_MAIL", "false") == "true" {
		// If the user is not verified, do not allow access.
		if user.VerifiedAt == nil {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "please verify your email address")
		}
	}

	// Send back JWT as a cookie.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		user.ID.String(),
		user.Username,
		c.IP(),
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"articpad-users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -2)),
			Issuer:    "articpad-api",
		},
	})
	signedToken, err := token.SignedString([]byte(config.GetString("SECRET", "MyRandomSecureSecret")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send response.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   signedToken,
	})
}

// Signs up a user and gives them a JWT.
func (h *AuthHandler) signUpUser(c *fiber.Ctx) error {
	// Create a struct so the request body can be mapped here.
	type registerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get request body.
	request := &registerRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Parse email.
	parsedEmail, err := mail.ParseAddress(request.Email)
	if err != nil || (err == nil && len(parsedEmail.Address) > 100) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidEmail)
	}

	// Check if username is valid.
	usernameValidator := validator.New(
		validator.MinLength(3, errors.New(consts.ErrUsernameLengthLessThan3)),
		validator.MaxLength(32, errors.New(consts.ErrUsernameLengthMoreThan32)),
		validator.ContainsOnly("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-_", errors.New(consts.ErrUsernameContainsInvalidCharacters)),
	)
	if err := usernameValidator.Validate(request.Username); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// Check if password is valid.
	similarity := 0.7
	passwordValidator := validator.New(
		validator.MinLength(8, errors.New(consts.ErrPasswordLengthLessThan8)),
		validator.MaxLength(64, errors.New(consts.ErrPasswordLengthMoreThan64)),
		validator.PasswordStrength(errors.New(consts.ErrPasswordStrength)),
		validator.Similarity([]string{request.Username, parsedEmail.Address}, &similarity, errors.New(consts.ErrPasswordSimilarity)),
	)
	if err := passwordValidator.Validate(request.Password); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// Check if user already exists.
	foundUser, err := h.userService.GetUserByEmail(customContext, parsedEmail.Address)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrEmailAlreadyExists)
	}

	// Check if username already exists.
	foundUser, err = h.userService.GetUserByUsername(customContext, request.Username)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrUsernameAlreadyExists)
	}

	// Hash password.
	// TODO: Look into using Argon2id instead of bcrypt and create a function for this for reusability.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	isAdmin := false
	// Check if this is the first user to set them as admin.
	if ok, _ := h.userService.IsFirstUser(customContext); ok {
		isAdmin = true
	}

	// Create user.
	user := &user.User{
		Username:          request.Username,
		Email:             parsedEmail.Address,
		Password:          string(hashedPassword),
		VerifiedAt:        nil,
		VerificationToken: uuid.New().String(),
		IsAdmin:           isAdmin,
	}

	err = h.userService.CreateUser(customContext, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if config.GetString("ENABLE_MAIL", "false") == "true" {
		// Send verification email.
		err := h.mailer.SendMail(parsedEmail.Address, "Verify your account", fmt.Sprintf("Please verify your account by clicking this link: <a href=\"%s\">%s</a>", config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/verify/"+user.VerificationToken, config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/verify/"+user.VerificationToken))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Your account was created but there was an error sending the verification email. If you don't receive an email, please request a new verification email. Error: "+err.Error())
		}
	}

	// Return result.
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"success": true,
		"message": "user has been created successfully." + func() string {
			if config.GetString("ENABLE_MAIL", "false") == "true" {
				return " please check your email to verify your account."
			}
			return ""
		}(),
	})
}

// Logs out a user.
func (h *AuthHandler) logOutUser(c *fiber.Ctx) error {
	// TODO: Invalidate JWT.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "user has been logged out successfully",
	})
}

// Refreshes a JWT.
func (h *AuthHandler) refreshToken(c *fiber.Ctx) error {
	// Create a struct for our custom JWT payload.
	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		UserIP string `json:"user_ip"`
		jwt.RegisteredClaims
	}

	// Get JWT data.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	// Create a new JWT.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		claims["uid"].(string),
		claims["user"].(string),
		claims["user_ip"].(string),
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"articpad-users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -2)),
			Issuer:    "articpad-api",
		},
	})

	// Sign the new JWT.
	signedToken, err := token.SignedString([]byte(config.GetString("SECRET", "MyRandomSecureSecret")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// TODO: Invalidate old JWT.

	// Send response.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   signedToken,
	})
}

// Gets the current logged in user.
func (h *AuthHandler) getMe(c *fiber.Ctx) error {
	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get JWT data.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	// Get user.
	user, err := h.userService.GetUser(customContext, claims["uid"].(uuid.UUID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send response.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"user":    user,
	})
}

// Resends a verification email to the user.
func (h *AuthHandler) resendVerificationEmail(c *fiber.Ctx) error {
	if config.GetString("ENABLE_MAIL", "false") == "false" {
		return fiber.NewError(fiber.StatusBadRequest, "email verification is disabled")
	}

	// Create a struct so the request body can be mapped here.
	type RequestPayload struct {
		Login string `json:"login"`
	}

	// Create a struct so the request body can be mapped here.
	request := new(RequestPayload)

	// Parse request body.
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch user by email or username.
	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusBadRequest, "user not found")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send verification email.
	err = h.mailer.SendMail(user.Email, "Verify your account", fmt.Sprintf("Please verify your account by clicking this link: <a href=\"%s\">%s</a>", config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/verify/"+user.VerificationToken, config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/verify/"+user.VerificationToken))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "A verification email has been sent to your email address",
	})
}

// Verifies a user's email and activates their account.
func (h *AuthHandler) verifyUser(c *fiber.Ctx) error {
	// Get verification token from URL.
	verificationToken := c.Params("token")

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch user by verification token.
	err := h.userService.VerifyUser(customContext, verificationToken)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusBadRequest, "invalid verification token")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Your email has been verified successfully",
	})
}

// Sends a password reset email to the user.
func (h *AuthHandler) forgotPassword(c *fiber.Ctx) error {
	if config.GetString("ENABLE_MAIL", "false") == "false" {
		return fiber.NewError(fiber.StatusBadRequest, "mail is not enabled, please contact the administrator to reset your password")
	}

	// Create a struct so the request body can be mapped here.
	type RequestPayload struct {
		Email string `json:"email"`
	}

	// Create a struct so the request body can be mapped here.
	request := new(RequestPayload)

	// Parse request body.
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch user by email
	user, err := h.userService.GetUserByEmail(customContext, request.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusBadRequest, "user not found")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Set new password reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 4)
	err = h.userService.SetPasswordResetToken(customContext, user.ID, token, expiresAt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send password reset email.
	err = h.mailer.SendMail(user.Email, "Reset your password", fmt.Sprintf("A password reset token has been generated for your account. If you did not request this, please ignore this email. Otherwise, here is your password reset token: <b>%s</b> <br> This token will expire in 4 hours.", token))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "A password reset email with an unique code valid for 4 hours has been sent to your email address",
	})
}

// Resets a user's password.
func (h *AuthHandler) resetPassword(c *fiber.Ctx) error {
	// Create a struct so the request body can be mapped here.
	type RequestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	// Create a struct so the request body can be mapped here.
	request := new(RequestPayload)

	// Parse request body.
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch user by email.
	user, err := h.userService.GetUserByEmail(customContext, request.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusBadRequest, "user not found")
	} else if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if token is valid.
	if user.PasswordResetToken != request.Token {
		return fiber.NewError(fiber.StatusBadRequest, "invalid password reset token")
	}

	// Check if token is expired.
	if user.PasswordResetExpiresAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "password reset token has expired")
	}

	// Hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Set new password.
	user.Password = string(hashedPassword)

	// Update user's password.
	err = h.userService.UpdateUser(customContext, user.ID, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Your password has been reset successfully",
	})
}

// Gets a single user.
func (h *AuthHandler) getUser(c *fiber.Ctx) error {
	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedUserID := c.Params("userID")

	// Validate parameter.
	parsedUserID, err := uuid.Parse(targetedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Get one user.
	user, err := h.userService.GetUser(customContext, parsedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return results.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// Updates a single user.
func (h *AuthHandler) updateUser(c *fiber.Ctx) error {
	// Initialize variables.
	user := &user.User{}
	targetedUserID := c.Locals("userID").(uuid.UUID)

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse request body.
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Update one user.
	err := h.userService.UpdateUser(customContext, targetedUserID, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "User has been updated successfully",
	})
}

// Deletes a single user.
func (h *AuthHandler) deleteUser(c *fiber.Ctx) error {
	// Initialize previous user ID.
	targetedUserID := c.Locals("userID").(uuid.UUID)

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete one user.
	err := h.userService.DeleteUser(customContext, targetedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return 204 No Content.
	return c.SendStatus(fiber.StatusNoContent)
}
