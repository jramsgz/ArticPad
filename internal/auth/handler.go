package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/internal/utils/consts"
	"github.com/jramsgz/articpad/internal/utils/templates"
	"github.com/jramsgz/articpad/pkg/apierror"
	"github.com/jramsgz/articpad/pkg/argon2id"
	"github.com/jramsgz/articpad/pkg/i18n"
	mailClient "github.com/jramsgz/articpad/pkg/mail"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userService user.UserService
	mailer      *mailClient.Mailer
	i18n        *i18n.I18n
}

type jwtClaims struct {
	UserID string `json:"uid"`
	User   string `json:"user"`
	UserIP string `json:"user_ip"`
	jwt.RegisteredClaims
}

// Creates a new authentication handler.
func NewAuthHandler(authRoute fiber.Router, us user.UserService, mail *mailClient.Mailer, i18n *i18n.I18n) {
	handler := &AuthHandler{
		userService: us,
		mailer:      mail,
		i18n:        i18n,
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
	authRoute.Get("/test", handler.test)                             // TODO
}

func (h *AuthHandler) test(c *fiber.Ctx) error {
	err := h.userService.DeleteUser(context.Background(), uuid.MustParse("ba5a7785-e31b-4159-be63-9658e67866fc"))
	return c.JSON(err)
}

// Signs in a user and gives them a JWT.
func (h *AuthHandler) signInUser(c *fiber.Ctx) error {
	type loginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &loginRequest{}
	if err := c.BodyParser(request); err != nil {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	langCode := h.getLangCode(c)

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil && (err == gorm.ErrRecordNotFound || err.Error() == consts.ErrDeletedRecord) {
		return apierror.NewApiError(fiber.StatusUnprocessableEntity, consts.ErrCodeAccountNotFound, h.i18n.T(langCode, "errors.account_not_found"))
	} else if err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	if ok, err := argon2id.ComparePasswordAndHash(request.Password, user.Password); err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	} else if !ok {
		return apierror.NewApiError(fiber.StatusUnauthorized, consts.ErrCodeInvalidCredentials, h.i18n.T(langCode, "errors.invalid_credentials"))
	}

	if config.GetString("ENABLE_MAIL") == "true" {
		if !user.VerifiedAt.Valid || user.VerifiedAt.Time.IsZero() || user.VerifiedAt.Time.After(time.Now()) {
			return apierror.NewApiError(fiber.StatusUnprocessableEntity, consts.ErrCodeEmailNotVerified, h.i18n.T(langCode, "errors.email_not_verified"))
		}
	}

	signedToken, err := newJWTToken(user.ID.String(), user.Username, c.IP())
	if err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
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
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	langCode := h.getLangCode(c)
	user := &user.User{
		Username:          request.Username,
		Email:             request.Email,
		Password:          request.Password,
		VerifiedAt:        sql.NullTime{Valid: false, Time: time.Time{}},
		VerificationToken: uuid.New().String(),
		Lang:              langCode,
	}

	err := h.userService.CreateUser(customContext, user)
	if err != nil {
		return consts.MapApiError(err, h.i18n, langCode)
	}

	if config.GetString("ENABLE_MAIL") == "true" {
		err := h.mailer.SendMail(templates.GetEmailVerificationEmail(h.i18n, user))
		if err != nil {
			return apierror.NewApiError(
				fiber.StatusInternalServerError, consts.ErrCodeCannotSendVerificationEmail,
				h.i18n.Ts(langCode, "errors.cannot_send_verification_email", "error", err.Error()),
			).ShowError()
		}
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"success": true,
		"message": h.i18n.T(langCode, "messages.account_created"),
	})
}

// Logs out a user.
func (h *AuthHandler) logOutUser(c *fiber.Ctx) error {
	// TODO: Invalidate JWT.
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": "You have been logged out successfully",
	})
}

// Refreshes a JWT.
func (h *AuthHandler) refreshToken(c *fiber.Ctx) error {
	jwtData := c.Locals("user").(*jwt.Token)
	claims := jwtData.Claims.(jwt.MapClaims)

	signedToken, err := newJWTToken(claims["uid"].(string), claims["user"].(string), claims["user_ip"].(string))
	if err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
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

	langCode := h.getLangCode(c)

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	user, err := h.userService.GetUser(customContext, uuid.MustParse(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierror.NewApiError(fiber.StatusUnauthorized, consts.ErrCodeAccountNotFound, h.i18n.T(langCode, "errors.account_not_found"))
		}
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"user":    user,
	})
}

// Resends a verification email to the user.
func (h *AuthHandler) resendVerificationEmail(c *fiber.Ctx) error {
	langCode := h.getLangCode(c)

	if config.GetString("ENABLE_MAIL") == "false" {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeMailNotEnabled, h.i18n.T(langCode, "errors.mail_not_enabled"))
	}

	type RequestPayload struct {
		Login string `json:"login"`
	}

	request := new(RequestPayload)
	err := c.BodyParser(request)
	if err != nil {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err.Error() == consts.ErrDeletedRecord {
			return apierror.NewApiError(fiber.StatusUnprocessableEntity, consts.ErrCodeAccountNotFound, h.i18n.T(langCode, "errors.account_not_found"))
		}
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	err = h.mailer.SendMail(templates.GetEmailVerificationEmail(h.i18n, user))
	if err != nil {
		return apierror.NewApiError(
			fiber.StatusInternalServerError, consts.ErrCodeCannotSendVerificationEmail,
			h.i18n.Ts(langCode, "errors.cannot_send_verification_email", "error", err.Error()),
		).ShowError()
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": h.i18n.T(langCode, "messages.verification_email_sent"),
	})
}

// Verifies a user's email and activates their account.
func (h *AuthHandler) verifyUser(c *fiber.Ctx) error {
	verificationToken := c.Params("token")
	langCode := h.getLangCode(c)

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := h.userService.VerifyUser(customContext, verificationToken)
	if err != nil && err == gorm.ErrRecordNotFound {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeInvalidVerificationToken, h.i18n.T(langCode, "errors.invalid_verification_token"))
	} else if err != nil {
		return consts.MapApiError(err, h.i18n, langCode)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": h.i18n.T(langCode, "messages.account_verified"),
	})
}

// Sends a password reset email to the user.
func (h *AuthHandler) forgotPassword(c *fiber.Ctx) error {
	langCode := h.getLangCode(c)

	if config.GetString("ENABLE_MAIL") == "false" {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeMailNotEnabled, h.i18n.T(langCode, "errors.mail_not_enabled_reset_password"))
	}

	type RequestPayload struct {
		Login string `json:"login"`
	}

	request := new(RequestPayload)
	err := c.BodyParser(request)
	if err != nil {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := h.userService.GetUserByEmailOrUsername(customContext, request.Login)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err.Error() == consts.ErrDeletedRecord {
			return apierror.NewApiError(fiber.StatusUnprocessableEntity, consts.ErrCodeAccountNotFound, h.i18n.T(langCode, "errors.account_not_found"))
		}
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 4)
	err = h.userService.SetPasswordResetToken(customContext, user.ID, token, expiresAt)
	if err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	err = h.mailer.SendMail(templates.GetPasswordResetEmail(h.i18n, user, token))
	if err != nil {
		return apierror.NewApiError(
			fiber.StatusInternalServerError, consts.ErrCodeCannotSendPasswordResetEmail,
			h.i18n.Ts(langCode, "errors.cannot_send_password_reset_email", "error", err.Error()),
		).ShowError()
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": h.i18n.T(langCode, "messages.password_reset_email_sent"),
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
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	langCode := h.getLangCode(c)

	err = h.userService.ResetPassword(customContext, request.Token, request.Password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeInvalidPasswordResetToken, h.i18n.T(langCode, "errors.invalid_password_reset_token"))
		}
		return consts.MapApiError(err, h.i18n, langCode)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"message": h.i18n.T(langCode, "messages.password_reset"),
	})
}

// Gets a single user.
func (h *AuthHandler) getUser(c *fiber.Ctx) error {
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetedUserID := c.Params("userID")

	parsedUserID, err := uuid.Parse(targetedUserID)
	if err != nil {
		return apierror.NewApiError(fiber.StatusBadRequest, consts.ErrCodeBadRequest, err.Error())
	}

	user, err := h.userService.GetUser(customContext, parsedUserID)
	if err != nil {
		return apierror.NewApiError(fiber.StatusInternalServerError, consts.ErrCodeUnknown, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *AuthHandler) getLangCode(c *fiber.Ctx) string {
	return h.i18n.ParseLanguage(c.Get("Accept-Language"))
}

func newJWTToken(userId string, username string, userIP string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		userId,
		username,
		userIP,
		jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"articpad-users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -2)),
			Issuer:    "articpad-api",
		},
	})
	return token.SignedString([]byte(config.GetString("SECRET")))
}
