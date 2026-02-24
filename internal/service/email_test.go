package service

import (
	"context"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/config"
)

func TestNewEmailService_Disabled(t *testing.T) {
	cfg := &config.Config{
		SMTPHost: "",
		SMTPPort: 587,
	}

	svc := NewEmailService(cfg)

	if svc.IsEnabled() {
		t.Error("IsEnabled() = true, want false when SMTP_HOST is empty")
	}
}

func TestNewEmailService_Enabled(t *testing.T) {
	cfg := &config.Config{
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		SMTPUser:     "user",
		SMTPPassword: "pass",
		SMTPFrom:     "noreply@example.com",
	}

	svc := NewEmailService(cfg)

	if !svc.IsEnabled() {
		t.Error("IsEnabled() = false, want true when SMTP_HOST is set")
	}
}

func TestSendEmail_DisabledMode(t *testing.T) {
	cfg := &config.Config{
		SMTPHost: "",
		SMTPPort: 587,
	}
	svc := NewEmailService(cfg)

	err := svc.SendEmail(context.Background(), "to@example.com", "Test Subject", "<p>Hello</p>")
	if err != nil {
		t.Errorf("SendEmail() error = %v, want nil in disabled mode", err)
	}
}

func TestSendEmail_EnabledWithInvalidHost(t *testing.T) {
	cfg := &config.Config{
		SMTPHost:     "invalid.host.that.does.not.exist.example.com",
		SMTPPort:     587,
		SMTPUser:     "user",
		SMTPPassword: "pass",
		SMTPFrom:     "noreply@example.com",
	}
	svc := NewEmailService(cfg)

	err := svc.SendEmail(context.Background(), "to@example.com", "Test", "<p>Hello</p>")
	if err == nil {
		t.Error("SendEmail() error = nil, want error when SMTP host is unreachable")
	}
}

func TestBuildMIMEMessage(t *testing.T) {
	msg := buildMIMEMessage("from@example.com", "to@example.com", "Test Subject", "<p>Hello</p>")

	checks := []string{
		"From: from@example.com\r\n",
		"To: to@example.com\r\n",
		"Subject: Test Subject\r\n",
		"MIME-Version: 1.0\r\n",
		"Content-Type: text/html; charset=\"UTF-8\"\r\n",
		"<p>Hello</p>",
	}
	for _, want := range checks {
		if !contains(msg, want) {
			t.Errorf("buildMIMEMessage() missing %q", want)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
