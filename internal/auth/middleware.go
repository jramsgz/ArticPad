package auth

import (
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/pkg/apierror"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Guards a specific endpoint in the API.
func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.GetString(config.Secret)),
		ErrorHandler: jwtError,
	})
}

// JWT error message.
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeInvalidJWT, "Missing or malformed JWT")
	}
	return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeInvalidJWT, "Invalid or expired JWT")
}

// Gets user data (their ID) from the JWT middleware. Should be executed after calling 'JWTMiddleware()'.
func GetDataFromJWT(c *fiber.Ctx) error {
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)
	parsedUserID, err := uuid.Parse(claims["uid"].(string))
	if err != nil {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeInvalidJWT, "Invalid or expired JWT")
	}

	c.Locals("currentUser", parsedUserID)
	return c.Next()
}
