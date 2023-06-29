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
	authRoute.Post("/forgot", handler.forgotPassword)                // TODO
	authRoute.Get("/reset", handler.resetPassword)                   // TODO
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
		return fiber.NewError(fiber.StatusUnprocessableEntity, "incorrect email, username or password")
	} else if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// If password is incorrect, do not allow access.
	if !h.checkPasswordHash(request.Password, user.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "incorrect email, username or password")
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
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid email")
	}

	// Check if username is valid.
	usernameValidator := validator.New(
		validator.MinLength(3, errors.New("username must be at least 3 characters")),
		validator.MaxLength(32, errors.New("username must be at most 32 characters")),
		validator.ContainsOnly("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-_", errors.New("username must only contain letters, numbers, dashes, underscores and dots")),
	)
	if err := usernameValidator.Validate(request.Username); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// Check if password is valid.
	similarity := 0.7
	passwordValidator := validator.New(
		validator.MinLength(8, errors.New("password must be at least 8 characters")),
		validator.MaxLength(64, errors.New("password must be at most 64 characters")),
		validator.PasswordStrength(nil),
		validator.Similarity([]string{request.Username, parsedEmail.Address}, &similarity, errors.New("password must not be too similar to username or email")),
	)
	if err := passwordValidator.Validate(request.Password); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	// Check if user already exists.
	foundUser, err := h.userService.GetUserByEmail(customContext, parsedEmail.Address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "this email is already in use")
	}

	// Check if username already exists.
	foundUser, err = h.userService.GetUserByUsername(customContext, request.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "username already exists")
	}

	// Hash password.
	// TODO: Look into using Argon2id instead of bcrypt.
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
	// Give form to our output response.
	type jwtResponse struct {
		UserID interface{} `json:"uid"`
		User   interface{} `json:"user"`
		Iss    interface{} `json:"iss"`
		Aud    interface{} `json:"aud"`
		Exp    interface{} `json:"exp"`
	}

	// Prepare our variables to be displayed.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	// Shape output response.
	jwtResp := &jwtResponse{
		UserID: claims["uid"],
		User:   claims["user"],
		Iss:    claims["iss"],
		Aud:    claims["aud"],
		Exp:    claims["exp"],
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Welcome to the private route",
		"jwtData": jwtResp,
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

	// Set new password reset token.
	// TODO

	// Send password reset email.
	err = h.mailer.SendMail(user.Email, "Reset your password", fmt.Sprintf("Please reset your password by clicking this link: <a href=\"%s\">%s</a>", config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/reset-password/"+user.PasswordResetToken, config.GetString("APP_URL", "http://localhost:8080")+"/api/v1/auth/reset-password/"+user.PasswordResetToken))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "A password reset email has been sent to your email address",
	})
}

// Resets a user's password.
func (h *AuthHandler) resetPassword(c *fiber.Ctx) error {
	// Create cancellable context.
	/*customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch request body.
	request := new(user.ResetPasswordRequest)
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body.
	err := h.validator.Struct(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Reset password.
	err = h.userService.ResetPassword(customContext, request.Token, request.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}*/

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
	targetedUserID, err := c.ParamsInt("userID")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please specify a valid user ID")
	}

	// Get one user.
	user, err := h.userService.GetUser(customContext, targetedUserID)
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
	targetedUserID := c.Locals("userID").(int)

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
	targetedUserID := c.Locals("userID").(int)

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
