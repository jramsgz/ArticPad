package templates

import (
	"bytes"
	"text/template"

	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/pkg/i18n"
	"github.com/jramsgz/articpad/pkg/mail"
)

// GetEmailVerificationEmail returns the email verification email message.
func GetEmailVerificationEmail(i18n *i18n.I18n, user *user.User) *mail.MailMessage {
	lang := i18n.ParseLanguage(user.Lang)
	t := buildTemplate(lang, "email_verification.html", map[string]string{
		"URL":     config.GetString("APP_URL", "http://localhost:8080") + "/verify/" + user.VerificationToken,
		"subject": "Email verification",
	})
	return &mail.MailMessage{
		To:          []string{user.Email},
		Subject:     "Email verification",
		ContentType: mail.ContentTypeTextHTML,
		Body:        t,
	}
}

// GetPasswordResetEmail returns the password reset email message.
func GetPasswordResetEmail(i18n *i18n.I18n, user *user.User, token string) *mail.MailMessage {
	lang := i18n.ParseLanguage(user.Lang)
	t := buildTemplate(lang, "password_reset.html", map[string]string{
		"URL":     config.GetString("APP_URL", "http://localhost:8080") + "/password-reset/" + token,
		"subject": "Reset your password",
	})
	return &mail.MailMessage{
		To:          []string{user.Email},
		Subject:     "Reset your password",
		ContentType: mail.ContentTypeTextHTML,
		Body:        t,
	}
}

// buildTemplate builds the template with the given language, template type and data.
func buildTemplate(lang, templateType string, data map[string]string) string {
	path := config.GetString("TEMPLATES_DIR", "templates")

	temp := template.Must(template.ParseGlob(path + "/mail/*.html"))

	var body bytes.Buffer
	if err := temp.ExecuteTemplate(&body, templateType, data); err != nil {
		panic(err)
	}
	return body.String()
}
