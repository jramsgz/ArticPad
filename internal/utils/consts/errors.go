package consts

// Here we store all the defined errors the backend can return
const (
	// ErrInvalidCredentials is returned when the user tries to login with invalid credentials
	ErrInvalidCredentials = "user, email or password is incorrect"
	// ErrUsernameLengthLessThan3 is returned when the user tries to register with a username that is less than 3 characters
	ErrUsernameLengthLessThan3 = "username must be at least 3 characters"
	// ErrUsernameLengthMoreThan32 is returned when the user tries to register with a username that is more than 32 characters
	ErrUsernameLengthMoreThan32 = "username must be at most 32 characters"
	// ErrUsernameContainsInvalidCharacters is returned when the user tries to register with a username that contains invalid characters
	ErrUsernameContainsInvalidCharacters = "username must only contain letters, numbers, dashes, underscores and dots"
	// ErrPasswordLengthLessThan8 is returned when the user tries to register with a password that is less than 8 characters
	ErrPasswordLengthLessThan8 = "password must be at least 8 characters"
	// ErrPasswordLengthMoreThan64 is returned when the user tries to register with a password that is more than 64 characters
	ErrPasswordLengthMoreThan64 = "password must be at most 64 characters"
	// ErrPasswordSimilarity is returned when the user tries to register with a password that is too similar to the username or email
	ErrPasswordSimilarity = "password must not be too similar to username or email"
	// ErrInvalidEmail is returned when the user tries to register with an invalid email
	ErrInvalidEmail = "invalid email"
	// ErrEmailAlreadyExists is returned when the user tries to register with an email that is already in use
	ErrEmailAlreadyExists = "this email is already in use"
	// ErrUsernameAlreadyExists is returned when the user tries to register with a username that is already in use
	ErrUsernameAlreadyExists = "username already exists"
	// ErrPasswordStrength is returned when the user tries to register with a password that is too weak
	ErrPasswordStrength = "password must contain at least one uppercase letter, one lowercase letter, one number and one special character"
	// ErrEmailNotVerified is returned when the user tries to login with an email that is not verified
	ErrEmailNotVerified = "please verify your email address"
)
