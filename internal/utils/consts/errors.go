package consts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jramsgz/articpad/pkg/apierror"
	"github.com/jramsgz/articpad/pkg/i18n"
)

// These are the defined errors returned using Go's errors package. They are used in the service and repository layers.
const (
	ErrUsernameLengthLessThan3           = "username must be at least 3 characters"
	ErrUsernameLengthMoreThan32          = "username must be at most 32 characters"
	ErrUsernameContainsInvalidCharacters = "username must only contain letters, numbers, dashes, underscores and dots"
	ErrPasswordLengthLessThan8           = "password must be at least 8 characters"
	ErrPasswordLengthMoreThan64          = "password must be at most 64 characters"
	ErrPasswordSimilarity                = "password must not be too similar to username or email"
	ErrInvalidEmail                      = "invalid email"
	ErrEmailAlreadyExists                = "this email is already in use"
	ErrUsernameAlreadyExists             = "username already exists"
	ErrPasswordStrength                  = "password must contain at least one uppercase letter, one lowercase letter, one number and one special character"
	ErrEmailAlreadyVerified              = "user is already verified"
	ErrPasswordResetTokenExpired         = "password reset token has expired"
	ErrDeletedRecord                     = "record has been deleted"
	ErrUsernameDeactivated               = "username has been deactivated"
	ErrEmailDeactivated                  = "email has been deactivated"
)

// These are the defined error codes returned by the API, most of them are related to errors defined by this package
const (
	ErrCodeUnknown                               = "unknown_error"
	ErrCodeBadRequest                            = "bad_request"
	ErrCodeAccountNotFound                       = "account_not_found"
	ErrCodeInvalidCredentials                    = "invalid_credentials"
	ErrCodeUsernameLengthLessThan3Code           = "username_length_less_than_3"
	ErrCodeUsernameLengthMoreThan32Code          = "username_length_more_than_32"
	ErrCodeUsernameContainsInvalidCharactersCode = "username_contains_invalid_characters"
	ErrCodePasswordLengthLessThan8Code           = "password_length_less_than_8"
	ErrCodePasswordLengthMoreThan64Code          = "password_length_more_than_64"
	ErrCodePasswordSimilarityCode                = "password_similarity"
	ErrCodeInvalidEmail                          = "invalid_email"
	ErrCodeEmailAlreadyExistsCode                = "email_already_exists"
	ErrCodeUsernameAlreadyExistsCode             = "username_already_exists"
	ErrCodePasswordStrengthCode                  = "password_strength"
	ErrCodeEmailNotVerified                      = "email_not_verified"
	ErrCodeEmailAlreadyVerifiedCode              = "email_already_verified"
	ErrCodeCannotSendVerificationEmail           = "cannot_send_verification_email"
	ErrCodeMailNotEnabled                        = "mail_not_enabled"
	ErrCodeCannotSendPasswordResetEmail          = "cannot_send_password_reset_email"
	ErrCodePasswordResetTokenExpired             = "password_reset_token_expired"
	ErrCodeInvalidJWT                            = "invalid_jwt"
	ErrCodeInvalidPasswordResetToken             = "invalid_password_reset_token"
	ErrCodeInvalidVerificationToken              = "invalid_verification_token"
	ErrCodeUsernameDeactivated                   = "username_deactivated"
	ErrCodeEmailDeactivated                      = "email_deactivated"
)

// appError is a struct that contains the data of an error returned by the API.
type appError struct {
	Status           int
	Code             string
	Message          string
	ForceShowMessage bool
}

// Map of most of the errors returned by underlying packages/layers.
var errorsMap = map[string]appError{
	ErrUsernameLengthLessThan3:           {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeUsernameLengthLessThan3Code, Message: "errors.username_too_short"},
	ErrUsernameLengthMoreThan32:          {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeUsernameLengthMoreThan32Code, Message: "errors.username_too_long"},
	ErrUsernameContainsInvalidCharacters: {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeUsernameContainsInvalidCharactersCode, Message: "errors.username_contains_invalid_characters"},
	ErrPasswordLengthLessThan8:           {Status: fiber.StatusUnprocessableEntity, Code: ErrCodePasswordLengthLessThan8Code, Message: "errors.password_too_short"},
	ErrPasswordLengthMoreThan64:          {Status: fiber.StatusUnprocessableEntity, Code: ErrCodePasswordLengthMoreThan64Code, Message: "errors.password_too_long"},
	ErrPasswordSimilarity:                {Status: fiber.StatusUnprocessableEntity, Code: ErrCodePasswordSimilarityCode, Message: "errors.password_too_similar"},
	ErrInvalidEmail:                      {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeInvalidEmail, Message: "errors.invalid_email"},
	ErrEmailAlreadyExists:                {Status: fiber.StatusConflict, Code: ErrCodeEmailAlreadyExistsCode, Message: "errors.email_already_exists"},
	ErrUsernameAlreadyExists:             {Status: fiber.StatusConflict, Code: ErrCodeUsernameAlreadyExistsCode, Message: "errors.username_already_exists"},
	ErrPasswordStrength:                  {Status: fiber.StatusUnprocessableEntity, Code: ErrCodePasswordStrengthCode, Message: "errors.password_not_strong_enough"},
	ErrEmailAlreadyVerified:              {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeEmailAlreadyVerifiedCode, Message: "errors.email_already_verified"},
	ErrPasswordResetTokenExpired:         {Status: fiber.StatusUnprocessableEntity, Code: ErrCodePasswordResetTokenExpired, Message: "errors.password_reset_token_expired"},
	ErrUsernameDeactivated:               {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeUsernameDeactivated, Message: "errors.username_deactivated"},
	ErrEmailDeactivated:                  {Status: fiber.StatusUnprocessableEntity, Code: ErrCodeEmailDeactivated, Message: "errors.email_deactivated"},
}

// MapApiError maps an error to the corresponding API error if possible.
// If the error is not mapped, it returns a generic error.
// This method is mostly used with errors returned by the service or repository layers.
// Errors returned by the API are created directly in the handlers.
func MapApiError(err error, i18n *i18n.I18n, langCode ...string) *apierror.Error {
	if appError, ok := errorsMap[err.Error()]; ok {
		if i18n != nil && len(langCode) > 0 {
			appError.Message = i18n.T(langCode[0], appError.Message)
		}
		return apierror.NewApiError(appError.Status, appError.Code, appError.Message)
	}

	return apierror.NewApiError(fiber.StatusInternalServerError, ErrCodeUnknown, err.Error())
}
