package mail

import (
	"testing"
)

func TestNewMailer(t *testing.T) {
	config := &MailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "testuser",
		Password: "testpassword",
		ForceTLS: true,
	}

	mailer, err := NewMailer(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mailer.Config != config {
		t.Errorf("expected Config to be %v, got %v", config, mailer.Config)
	}
	if mailer.Client == nil {
		t.Error("expected Client to be initialized")
	}
}
