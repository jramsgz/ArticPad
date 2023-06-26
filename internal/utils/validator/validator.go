package validator

import (
	"errors"
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
