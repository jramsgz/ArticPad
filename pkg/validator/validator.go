package validator

import (
	"errors"

	"github.com/jramsgz/articpad/internal/utils/consts"
)

// ValidateFunc defines a function to validate
type ValidateFunc func(text string) error

// ValidateConditionFunc defines a function to validate with condition
type ValidateConditionFunc func(text string, condition bool) error

// Validator represents set of text validators
type Validator []ValidateFunc

// New return new instance of Validator
func New(vfunc ...ValidateFunc) *Validator {
	v := Validator{}
	v = append(v, vfunc...)
	return &v
}

// Validate the text
func (v *Validator) Validate(text string) error {
	if len(*v) == 0 {
		return errors.New("Validator must contains at least 1 validator function")
	}
	for _, textValidator := range *v {
		err := textValidator(text)
		if err != nil {
			return err
		}
	}
	return nil
}

// Default validator functions
var (
	// Username validator
	DefaultUsernameValidator = New(
		MinLength(3, consts.ErrUsernameLengthLessThan3),
		MaxLength(32, consts.ErrUsernameLengthMoreThan32),
		ContainsOnly("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-_", consts.ErrUsernameContainsInvalidCharacters),
	)
	// Password validator
	DefaultPasswordValidator = func(similarityAttributes []string) *Validator {
		similarity := 0.7
		v := New(
			MinLength(8, consts.ErrPasswordLengthLessThan8),
			MaxLength(64, consts.ErrPasswordLengthMoreThan64),
			Regex(`^\P{Ll}*\p{Ll}[\s\S]*$`, consts.ErrPasswordStrength),
			Regex(`^\P{Lu}*\p{Lu}[\s\S]*$`, consts.ErrPasswordStrength),
			Regex(`^\P{N}*\p{N}[\s\S]*$`, consts.ErrPasswordStrength),
			Regex(`^[\p{L}\p{N}]*[^\p{L}\p{N}][\s\S]*$`, consts.ErrPasswordStrength),
			Similarity(similarityAttributes, &similarity, consts.ErrPasswordSimilarity),
		)
		return v
	}
)
