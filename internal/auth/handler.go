package auth

import (
	"context"
	"database/sql"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/internal/utils/i18n"
	"github.com/jramsgz/articpad/internal/utils/templates"
	"github.com/jramsgz/articpad/pkg/argon2id"
	mailClient "github.com/jramsgz/articpad/pkg/mail"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService user.UserService
	mailer      *mailClient.Mailer
}

// Creates a new authentication handler.
func NewAuthHandler(authRoute fiber.Router, us user.UserService, mail *mailClient.Mailer) {
	handler := &AuthHandler{
		userService: us,
		mailer:      mail,
	}

	authRoute.Post("/login", handler.signInUser)
	authRoute.Post("/register", handler.signUpUser)
	authRoute.Post("/logout", JWTMiddleware(), handler.logOutUser) // TODO
	authRoute.Post("/resend", handler.resendVerificationEmail)
	authRoute.Get("/verify/:token", handler.verifyUser)
	authRoute.Post("/forgot", handler.forgotPassword)
	authRoute.Post("/reset", handler.resetPassword)
	authRoute.Get("/refresh", JWTMiddleware(), handler.refreshToken) // TODO
	authRoute.Get("/me", JWTMiddleware(), handler.getMe)             // TODO
}

// Signs in a user and gives them a JWT.
func (h *AuthHandler) signInUser(c *fiber.Ctx) error {
	type loginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		UserIP string `json:"user_ip"`
		jwt.RegisteredClaims
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &loginRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidCredentials)
	} else if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := argon2id.ComparePasswordAndHash(request.Password, user.Password); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if !ok {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidCredentials)
	}

	if config.GetString("ENABLE_MAIL", "false") == "true" {
		if !user.VerifiedAt.Valid || user.VerifiedAt.Time.IsZero() || user.VerifiedAt.Time.Before(time.Now()) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrEmailNotVerified)
		}
	}

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

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   signedToken,
		"user":    user,
	})
}

// Signs up a user and gives them a JWT.
func (h *AuthHandler) signUpUser(c *fiber.Ctx) error {
	type registerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &registerRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	parsedEmail, err := mail.ParseAddress(request.Email)
	if err != nil || (err == nil && len(parsedEmail.Address) > 100) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, consts.ErrInvalidEmail)
	}

	isAdmin := false
	if ok, _ := h.userService.IsFirstUser(customContext); ok {
		isAdmin = true
	}

	user := &user.User{
		Username:          request.Username,
		Email:             parsedEmail.Address,
		Password:          request.Password,
		VerifiedAt:        sql.NullTime{Valid: false, Time: time.Time{}},
		VerificationToken: uuid.New().String(),
		IsAdmin:           isAdmin,
		Lang:              i18n.ParseLanguageHeader(c.Get("Accept-Language")).String(),
	}

	err = h.userService.CreateUser(customContext, user)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return fiberErr
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if config.GetString("ENABLE_MAIL", "false") == "true" {
		err := h.mailer.SendMail(templates.GetEmailVerificationEmail(user))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Your account was created but there was an error sending the verification email. If you don't receive an email, please request a new verification email. Error: "+err.Error())
		}
	}

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
	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		UserIP string `json:"user_ip"`
		jwt.RegisteredClaims
	}

	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

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

	signedToken, err := token.SignedString([]byte(config.GetString("SECRET", "MyRandomSecureSecret")))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// TODO: Invalidate old JWT.

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   signedToken,
	})
}

// Gets the current logged in user.
func (h *AuthHandler) getMe(c *fiber.Ctx) error {
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	user, err := h.userService.GetUser(customContext, claims["uid"].(uuid.UUID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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

	type RequestPayload struct {
		Login string `json:"login"`
	}

	request := new(RequestPayload)
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err != gorm.ErrRecordNotFound {
		err = h.mailer.SendMail(templates.GetEmailVerificationEmail(user))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "If your email exists in our database, a verification email has been sent to it.",
	})
}

// Verifies a user's email and activates their account.
func (h *AuthHandler) verifyUser(c *fiber.Ctx) error {
	verificationToken := c.Params("token")

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := h.userService.VerifyUser(customContext, verificationToken)
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusBadRequest, "invalid verification token")
	} else if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return fiberErr
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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

	type RequestPayload struct {
		Login string `json:"login"`
	}

	request := new(RequestPayload)
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err != gorm.ErrRecordNotFound {
		token := uuid.New().String()
		expiresAt := time.Now().Add(time.Hour * 4)
		err = h.userService.SetPasswordResetToken(customContext, user.ID, token, expiresAt)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = h.mailer.SendMail(templates.GetPasswordResetEmail(user, token))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "If your email address or username exists in our database, you will receive a password recovery link valid for 4 hours at your email address in a few minutes",
	})
}

// Resets a user's password.
func (h *AuthHandler) resetPassword(c *fiber.Ctx) error {
	type RequestPayload struct {
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	request := new(RequestPayload)
	err := c.BodyParser(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = h.userService.ResetPassword(customContext, request.Token, request.Password)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return fiberErr
		} else if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "invalid password reset token")
		} else if err.Error() == "token has expired" {
			return fiber.NewError(fiber.StatusBadRequest, "password reset token has expired")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Your password has been reset successfully",
	})
}

// Gets a single user.
func (h *AuthHandler) getUser(c *fiber.Ctx) error {
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetedUserID := c.Params("userID")

	parsedUserID, err := uuid.Parse(targetedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.userService.GetUser(customContext, parsedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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
