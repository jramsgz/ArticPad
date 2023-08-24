package validator

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// MinLength returns a ValidateFunc that check if text length is not lower that "length"
func MinLength(length int, customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		if len(text) < length {
			if customError != nil {
				return customError
			}
			return fmt.Errorf("length must be not lower that %d chars", length)
		}
		return nil
	})
}

// MaxLength returns a ValidateFunc that check if text length is not greater that "length"
func MaxLength(length int, customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		if len(text) > length {
			if customError != nil {
				return customError
			}
			return fmt.Errorf("length must be not greater that %d chars", length)
		}
		return nil
	})
}

// ContainsOnly returns a ValidateFunc that check if text contains only selected chars
func ContainsOnly(chars string, customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		for _, char := range text {
			idx := strings.IndexFunc(chars, func(r rune) bool {
				return r == char
			})
			if idx == -1 {
				if customError != nil {
					return customError
				}
				return fmt.Errorf("must contain only %s", chars)
			}
		}
		return nil
	})
}

// ContainsAtLeast returns a ValidateFunc that count occurrences of a chars and compares it with required value
func ContainsAtLeast(chars string, occurrences int, customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		cnt := 0
		for _, char := range strings.Split(chars, "") {
			cnt += strings.Count(text, char)
		}
		if cnt < occurrences {
			if customError != nil {
				return customError
			}
			return fmt.Errorf("must contain at least %d chars from %s", occurrences, chars)
		}
		return nil
	})
}

// PasswordStrength returns a ValidateFunc that check if text contains at least one lowercase, uppercase, digit and special char
func PasswordStrength(customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		var (
			hasUpper   = false
			hasLower   = false
			hasNumber  = false
			hasSpecial = false
		)

		for _, c := range text {
			switch {
			case unicode.IsNumber(c):
				hasNumber = true
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsLower(c):
				hasLower = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSpecial = true
			}
		}

		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			if customError != nil {
				return customError
			}
		}

		if !hasUpper {
			return errors.New("password must contain at least one uppercase letter")
		}
		if !hasLower {
			return errors.New("password must contain at least one lowercase letter")
		}
		if !hasNumber {
			return errors.New("password must contain at least one number")
		}
		if !hasSpecial {
			return errors.New("password must contain at least one special character")
		}

		return nil
	})
}

// AlphaNumeric returns a ValidateFunc that check if text contains only letters and digits
func AlphaNumeric(customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		for _, char := range text {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
				if customError != nil {
					return customError
				}
				return errors.New("must contain only letters and digits")
			}
		}
		return nil
	})
}

// Similarity returns ValidateFunc that validate whether the text is sufficiently different from the attributes
func Similarity(attributes []string, maxSimilarity *float64, customError error) ValidateFunc {
	if maxSimilarity == nil {
		maxSimilarity = new(float64)
		*maxSimilarity = 0.7
	}
	return ValidateFunc(func(text string) error {
		for idx := range attributes {
			if ratio(text, attributes[idx]) > *maxSimilarity {
				if customError != nil {
					return customError
				}
				return fmt.Errorf("too similar to the %s", attributes[idx])
			}
		}
		return nil
	})
}

// Return a measure of the sequences' similarity (float in [0,1]).
//
// Where T is the total number of elements in both sequences, and M is the number of matches, this is 2.0*M / T.
// Note that this is 1 if the sequences are identical, and 0 if they have nothing in common.
func ratio(a, b string) float64 {
	t := float64(len(a) + len(b))
	m := 0.0
	sa := strings.Split(a, "")
	sb := strings.Split(b, "")
	for idx := range sb {
		aidx := strings.Index(a, sb[idx])
		if aidx > -1 {
			sa = append(sa[:aidx], sa[aidx+1:]...)
			a = strings.Join(sa, "")
			m = m + 1.0
		}
	}
	return 2.0 * m / t
}

// StartsWith returns ValidateFunc that validate whether the text is starts with one of letter
func StartsWith(letters string, customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		firstLetter := []rune(text)[0]
		idx := strings.IndexFunc(letters, func(r rune) bool {
			return r == firstLetter
		})
		if idx == -1 {
			if customError != nil {
				return customError
			}
			return fmt.Errorf("must start with one of: %s", letters)
		}
		return nil
	})
}

// Unique returns ValidateFunc that validate whether the text has only unique chars
func Unique(customError error) ValidateFunc {
	return ValidateFunc(func(text string) error {
		runes := []rune(text)
		for idx := range runes {
			if strings.LastIndexFunc(text, func(r rune) bool {
				return r == runes[idx]
			}) != idx {
				if customError != nil {
					return customError
				}
				return errors.New("must contain unique chars")
			}
		}
		return nil
	})
}
