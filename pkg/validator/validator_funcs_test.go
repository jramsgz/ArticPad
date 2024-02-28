package validator

import (
	"errors"
	"testing"
)

func TestPasswordStrength(t *testing.T) {
	// Test case 1: Valid password
	err := PasswordStrength(nil)("Password1!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Missing uppercase letter
	err = PasswordStrength(nil)("password1!")
	expectedErr := errors.New("password must contain at least one uppercase letter")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 3: Missing lowercase letter
	err = PasswordStrength(nil)("PASSWORD1!")
	expectedErr = errors.New("password must contain at least one lowercase letter")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 4: Missing number
	err = PasswordStrength(nil)("Password!")
	expectedErr = errors.New("password must contain at least one number")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 5: Missing special character
	err = PasswordStrength(nil)("Password1")
	expectedErr = errors.New("password must contain at least one special character")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 6: Custom error
	customErr := errors.New("custom error")
	err = PasswordStrength(customErr)("password1!")
	if err != customErr {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestMinLength(t *testing.T) {
	// Test case 1: Valid length
	err := MinLength(3, nil)("abcd")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Invalid length
	expectedErr := errors.New("text length must be at least 5")
	err = MinLength(5, expectedErr)("abcd")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 3: Custom error
	customErr := errors.New("custom error")
	err = MinLength(5, customErr)("abcd")
	if err != customErr {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestMaxLength(t *testing.T) {
	// Test case 1: Valid length
	err := MaxLength(5, nil)("abcd")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Invalid length
	expectedErr := errors.New("text length must be at most 3")
	err = MaxLength(3, expectedErr)("abcd")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 3: Custom error
	customErr := errors.New("custom error")
	err = MaxLength(3, customErr)("abcd")
	if err != customErr {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestContainsOnly(t *testing.T) {
	// Test case 1: Valid characters
	err := ContainsOnly("abcd", nil)("abcd")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Invalid characters
	expectedErr := errors.New("text contains invalid characters")
	err = ContainsOnly("abcd", expectedErr)("abcde")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 3: Custom error
	customErr := errors.New("custom error")
	err = ContainsOnly("abcd", customErr)("abcde")
	if err != customErr {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestRegex(t *testing.T) {
	// Test case 1: Valid regex
	err := Regex(`^\d+$`, nil)("123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test case 2: Invalid regex
	expectedErr := errors.New("text does not match regex")
	err = Regex(`^\d+$`, expectedErr)("abc")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 3: Custom error
	customErr := errors.New("custom error")
	err = Regex(`^\d+$`, customErr)("abc")
	if err != customErr {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSimilarity(t *testing.T) {
	// Test case 1: Similarity
	similarity := 0.7
	expectedErr := errors.New("text is too similar to other text")
	err := Similarity([]string{"abc"}, &similarity, expectedErr)("abc")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("expected error %q, got %v", expectedErr, err)
	}

	// Test case 2: Not similar
	err = Similarity([]string{"abc"}, &similarity, nil)("def")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
