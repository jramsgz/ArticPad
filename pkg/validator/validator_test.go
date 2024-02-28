package validator

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := New()
	if len(*v) != 0 {
		t.Error("expected empty validator")
	}

	v = New(MinLength(3, nil))
	if len(*v) != 1 {
		t.Error("expected 1 validator")
	}

	err := v.Validate("abcd")
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
