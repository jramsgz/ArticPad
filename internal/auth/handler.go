package auth

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/mailer"
	"github.com/jramsgz/articpad/internal/user"
	"golang.org/x/crypto/bcrypt"
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
	authRoute.Post("/logout", JWTMiddleware(), handler.logOutUser)   // TODO
	authRoute.Get("/verify", handler.verifyUser)                     // TODO
	authRoute.Post("/forgot", handler.forgotPassword)                // TODO
	authRoute.Get("/reset", handler.resetPassword)                   // TODO
	authRoute.Get("/refresh", JWTMiddleware(), handler.refreshToken) // TODO
	authRoute.Get("/me", JWTMiddleware(), handler.getMe)             // TODO
}

// CheckPasswordHash compare password with hash
func (h *AuthHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Checks if an email is valid.
func (h *AuthHandler) validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Signs in a user and gives them a JWT.
func (h *AuthHandler) signInUser(c *fiber.Ctx) error {
	// Create a struct so the request body can be mapped here.
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Create a struct for our custom JWT payload.
	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		jwt.StandardClaims
	}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get request body.
	request := &loginRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Check if email is valid.
	if h.validEmail(request.Email) == false {
		return fiber.NewError(fiber.StatusBadRequest, "Please specify a valid email!")
	}

	// Get user from database.
	user, err := h.userService.GetUserByEmail(customContext, request.Email)
	if err != nil {
		// If the user does not exist, do not allow access.
		if err.Error() == "record not found" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Incorrect email or password!")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// If both email and password are incorrect, do not allow access.
	if h.checkPasswordHash(request.Password, user.Password) == false {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Incorrect email or password!")
	}

	if config.GetString("SEND_VERIFY_EMAIL", "false") == "true" {
		// If the user is not verified, do not allow access.
		if user.Verified == false {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "Please verify your account!")
		}
	}

	// Send back JWT as a cookie.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		user.ID.String(),
		user.Username,
		jwt.StandardClaims{
			Audience:  "articpad-users",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
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

	// Check if email is valid.
	if !h.validEmail(request.Email) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid email!")
	}

	// Check if username is valid.
	if len(request.Username) < 3 || len(request.Username) > 32 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Username must be between 3 and 32 characters!")
	}

	// Check if password is valid.
	if len(request.Password) < 8 || len(request.Password) > 64 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Password must be between 8 and 64 characters!")
	}

	// Check if user already exists.
	foundUser, err := h.userService.GetUserByEmail(customContext, request.Email)
	if err != nil {
		if err.Error() != "record not found" {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "This email is already in use!")
	}

	// Check if username already exists.
	foundUser, err = h.userService.GetUserByUsername(customContext, request.Username)
	if err != nil {
		if err.Error() != "record not found" {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	if foundUser != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Username already exists!")
	}

	// Hash password.
	// TODO: Look into using Argon2id instead of bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Create user.
	user := &user.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Verified: false,
		Admin:    false,
	}

	err = h.userService.CreateUser(customContext, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if config.GetString("SEND_VERIFY_EMAIL", "false") == "true" {
		// Send verification email.
		h.mailer.SendMail(request.Email, "Verify your account", fmt.Sprintf("Please verify your account by clicking <a href=\"%s\">here</a>.", config.GetString("APP_URL", "http://localhost:8080")))
	}

	// Return result.
	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"success": true,
		"message": "User has been created successfully! Please check your email to verify your account.",
	})
}

// Logs out a user.
func (h *AuthHandler) logOutUser(c *fiber.Ctx) error {
	// TODO: Invalidate JWT.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "User has been logged out successfully!",
	})
}

// Refreshes a JWT.
func (h *AuthHandler) refreshToken(c *fiber.Ctx) error {
	// Create a struct for our custom JWT payload.
	type jwtClaims struct {
		UserID string `json:"uid"`
		User   string `json:"user"`
		jwt.StandardClaims
	}

	// Get JWT data.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	// Create a new JWT.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		claims["uid"].(string),
		claims["user"].(string),
		jwt.StandardClaims{
			Audience:  "articpad-users",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
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
		"message": "Welcome to the private route!",
		"jwtData": jwtResp,
	})
}

// Verifies a user's email and activates their account.
func (h *AuthHandler) verifyUser(c *fiber.Ctx) error {
	// Create cancellable context.
	/*customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	verificationToken := c.Params("token")

	// Verify email.
	err := h.userService.VerifyEmail(customContext, verificationToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}*/

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "Your email has been verified successfully!",
	})
}

// Sends a password reset email to the user.
func (h *AuthHandler) forgotPassword(c *fiber.Ctx) error {
	// Create cancellable context.
	/*customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch request body.
	request := new(user.ForgotPasswordRequest)
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body.
	err := h.validator.Struct(request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Get user.
	user, err := h.userService.GetUserByEmail(customContext, request.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Generate a password reset token.
	passwordResetToken, err := h.userService.GeneratePasswordResetToken(customContext, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send password reset email.
	err = h.mailer.SendPasswordResetEmail(user.Email, passwordResetToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}*/

	// Return result.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "A password reset email has been sent to your email address!",
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
		"message": "Your password has been reset successfully!",
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
		return fiber.NewError(fiber.StatusBadRequest, "Please specify a valid user ID!")
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
		"message": "User has been updated successfully!",
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
