package templates

import (
	"bytes"
	"text/template"

	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/user"
	"github.com/jramsgz/articpad/internal/utils/i18n"
	"github.com/jramsgz/articpad/pkg/mail"
	"golang.org/x/text/language"
)

// Available languages.
const (
	English = "en"
	Spanish = "es"
)

// GetEmailVerificationEmail returns the email verification email message.
func GetEmailVerificationEmail(user *user.User) *mail.MailMessage {
	lang := i18n.ParseLanguage(user.Lang)
	t := buildTemplate(lang, "email_verification.html", map[string]string{
		"URL":     config.GetString("APP_URL", "http://localhost:8080") + "/password-reset/" + user.VerificationToken,
		"subject": "Email verification",
	})
	return &mail.MailMessage{
		To:      []string{user.Email},
		Subject: "Email verification",
		Body:    t,
	}
}

// GetPasswordResetEmail returns the password reset email message.
func GetPasswordResetEmail(user *user.User, token string) *mail.MailMessage {
	lang := i18n.ParseLanguage(user.Lang)
	t := buildTemplate(lang, "password_reset.html", map[string]string{
		"token":   token,
		"URL":     config.GetString("APP_URL", "http://localhost:3000"),
		"subject": "Reset your password",
	})
	return &mail.MailMessage{
		To:      []string{user.Email},
		Subject: "Reset your password",
		/* A password reset token has been generated for your account. If you did not request this, please ignore this email and do not share this token with anyone. If you want to reset your password, please click this link: <a href=\"%s\">%s</a>", config.GetString("APP_URL", "http://localhost:8080")+"/password-reset/"+token, config.GetString("APP_URL", "http://localhost:8080")+"/password-reset/"+token)) */
		Body: t,
	}
}

// buildTemplate builds the template with the given language, template type and data.
func buildTemplate(lang language.Tag, templateType string, data map[string]string) string {
	path := config.GetString("TEMPLATES_DIR", "templates")
	langPath := parseLanguagePath(lang)

	temp := template.Must(template.ParseGlob(path + "/mail/" + langPath + "/*.html"))

	var body bytes.Buffer
	if err := temp.ExecuteTemplate(&body, templateType, data); err != nil {
		panic(err)
	}
	return body.String()
}

// ParseLanguage parses the language string and returns the language code.
func parseLanguagePath(lang language.Tag) string {
	switch lang {
	case language.English:
		return English
	case language.Spanish:
		return Spanish
	default:
		return English
	}
}
