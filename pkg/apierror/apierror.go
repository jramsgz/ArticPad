package apierror

import "github.com/gofiber/fiber/v2/utils"

// Error represents an error that occurred while handling a request.
type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Show    bool   `json:"show"`
}

// Error makes it compatible with the `error` interface.
func (e *Error) Error() string {
	return e.Message
}

// ShowError sets the Show field to true meaning that the error should be shown to the user
// even if the application is running in production mode.
func (e *Error) ShowError() *Error {
	e.Show = true
	return e
}

// NewApiError creates a new Error instance with an optional message
func NewApiError(status int, code string, message ...string) *Error {
	err := &Error{
		Status:  status,
		Code:    code,
		Message: utils.StatusMessage(status),
		Show:    false,
	}
	if len(message) > 0 {
		err.Message = message[0]
	}
	return err
}
