package apierror

import (
	"testing"
)

func TestNewApiError(t *testing.T) {
	status := 404
	code := "not_found"
	message := "Resource not found"

	err := NewApiError(status, code, message)

	if err.Status != status {
		t.Errorf("Expected status %d, but got %d", status, err.Status)
	}

	if err.Code != code {
		t.Errorf("Expected code %s, but got %s", code, err.Code)
	}

	if err.Message != message {
		t.Errorf("Expected message %s, but got %s", message, err.Message)
	}

	if err.Show {
		t.Errorf("Expected show to be false, but got true")
	}
}

func TestError(t *testing.T) {
	status := 404
	code := "not_found"
	message := "Resource not found"

	err := NewApiError(status, code, message)

	if err.Error() != message {
		t.Errorf("Expected error message %s, but got %s", message, err.Error())
	}
}

func TestShowError(t *testing.T) {
	status := 404
	code := "not_found"
	message := "Resource not found"

	err := NewApiError(status, code, message)

	if err.Show {
		t.Errorf("Expected show to be false, but got true")
	}

	err.ShowError()

	if !err.Show {
		t.Errorf("Expected show to be true, but got false")
	}
}
