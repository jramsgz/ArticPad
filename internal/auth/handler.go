package auth

import (
	"net/mail"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jramsgz/articpad/config"
	"golang.org/x/crypto/bcrypt"
)

// Create an authentication handler. Leave this empty, as we have no domains nor use-cases.
type AuthHandler struct{}

// Creates a new authentication handler.
func NewAuthHandler(authRoute fiber.Router) {
	handler := &AuthHandler{}

	// Declare routing for specific routes.
	authRoute.Post("/login", handler.signInUser)
	authRoute.Get("/refresh", JWTMiddleware(), handler.refreshToken)
	authRoute.Get("/private", JWTMiddleware(), handler.privateRoute)
}

// CheckPasswordHash compare password with hash
func (h *AuthHandler) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Checks if an email is valid.
func (h *AuthHandler) valid(email string) bool {
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

	// Get request body.
	request := &loginRequest{}
	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// If both username and password are incorrect, do not allow access.
	if request.Email != os.Getenv("API_USERNAME") || request.Password != os.Getenv("API_PASSWORD") {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong username or password!")
	}

	// Send back JWT as a cookie.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		"123",
		"username",
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

	// Send response.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"token":   signedToken,
	})
}

// A single private route that only logged in users can access.
func (h *AuthHandler) privateRoute(c *fiber.Ctx) error {
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
