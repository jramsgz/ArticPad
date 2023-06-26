package auth

import (
	"context"
	"net/mail"
	"strconv"

	"github.com/jramsgz/articpad/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Guards a specific endpoint in the API.
func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.GetString("SECRET", "MyRandomSecureSecret")),
		ErrorHandler: jwtError,
	})
}

// JWT error message.
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return fiber.NewError(fiber.StatusBadRequest, "Invalid or expired JWT")
}

// Gets user data (their ID) from the JWT middleware. Should be executed after calling 'JWTMiddleware()'.
func GetDataFromJWT(c *fiber.Ctx) error {
	// Get userID from the previous route.
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)
	parsedUserID := claims["uid"].(string)
	userID, err := strconv.Atoi(parsedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Go to next.
	c.Locals("currentUser", userID)
	return c.Next()
}

// If user does not exist, do not allow one to access the API.
func (h *AuthHandler) checkIfUserExistsMiddleware(c *fiber.Ctx) error {
	// Create a new customized context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedUserEmail := c.Params("email")
	parsedEmail, err := mail.ParseAddress(targetedUserEmail)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email address")
	}

	// Check if user exists.
	searchedUser, err := h.userService.GetUserByEmail(customContext, parsedEmail.Address)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if searchedUser == nil {
		return fiber.NewError(fiber.StatusBadRequest, "There is no user with this email!")
	}

	// Store in locals for further processing in the real handler.
	c.Locals("userID", searchedUser.ID)
	return c.Next()
}
