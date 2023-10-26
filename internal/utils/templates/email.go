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
	t := buildTemplate("email_verification.html", map[string]string{
		"URL":        config.GetString("APP_URL") + "/verify/" + user.VerificationToken,
		"Subject":    i18n.T(lang, "email.verification.subject"),
		"Header":     i18n.T(lang, "email.verification.header"),
		"LogoURL":    config.GetString("APP_URL") + "/assets/logo_vertical.png",
		"Title":      i18n.T(lang, "email.verification.title"),
		"Content":    i18n.T(lang, "email.verification.content"),
		"Button":     i18n.T(lang, "email.verification.button"),
		"ButtonLink": i18n.T(lang, "email.button_link"),
		"Footer":     i18n.T(lang, "email.footer"),
	})
	return &mail.MailMessage{
		To:          []string{user.Email},
		Subject:     i18n.T(lang, "email.verification.subject"),
		ContentType: mail.ContentTypeTextHTML,
		Body:        t,
	}
}

// GetPasswordResetEmail returns the password reset email message.
func GetPasswordResetEmail(i18n *i18n.I18n, user *user.User, token string) *mail.MailMessage {
	lang := i18n.ParseLanguage(user.Lang)
	t := buildTemplate("password_reset.html", map[string]string{
		"URL":        config.GetString("APP_URL") + "/password-reset/" + token,
		"Subject":    i18n.T(lang, "email.password_reset.subject"),
		"Header":     i18n.T(lang, "email.password_reset.header"),
		"LogoURL":    config.GetString("APP_URL") + "/assets/logo_vertical.png",
		"Title":      i18n.T(lang, "email.password_reset.title"),
		"Content":    i18n.T(lang, "email.password_reset.content"),
		"Button":     i18n.T(lang, "email.password_reset.button"),
		"ButtonLink": i18n.T(lang, "email.button_link"),
		"Footer":     i18n.T(lang, "email.footer"),
	})
	return &mail.MailMessage{
		To:          []string{user.Email},
		Subject:     i18n.T(lang, "email.password_reset.subject"),
		ContentType: mail.ContentTypeTextHTML,
		Body:        t,
	}
}

// buildTemplate builds the template with the given language, template type and data.
func buildTemplate(templateType string, data map[string]string) string {
	path := config.GetString("TEMPLATES_DIR")

	temp := template.Must(template.ParseGlob(path + "/mail/*.html"))

	var body bytes.Buffer
	if err := temp.ExecuteTemplate(&body, templateType, data); err != nil {
		panic(err)
	}
	return body.String()
}
