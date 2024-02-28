package mail

import (
	"context"

	"github.com/wneessen/go-mail"
)

// MailConfig represents the configuration for the mail server.
type MailConfig struct {
	From     string
	Username string
	Password string
	Port     int
	Host     string
	ForceTLS bool
}

// Mailer represents a mailer client.
type Mailer struct {
	Config *MailConfig
	Client *mail.Client
}

// Mail represents a mail message.
type MailMessage struct {
	To          []string
	Subject     string
	ContentType ContentType
	Body        string
}

// ContentType represents the content type of a mail message.
type ContentType string

const (
	// ContentTypeTextPlain represents a plain text content type.
	ContentTypeTextPlain ContentType = "text/plain"
	// ContentTypeTextHTML represents a HTML content type.
	ContentTypeTextHTML ContentType = "text/html"
)

// NewMailer creates a new mailer client and connects to the mail server.
func NewMailer(config *MailConfig) (*Mailer, error) {
	TLSPolicy := mail.TLSOpportunistic
	if config.ForceTLS {
		TLSPolicy = mail.TLSMandatory
	}

	c, err := mail.NewClient(config.Host, mail.WithPort(config.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPolicy(TLSPolicy),
		mail.WithUsername(config.Username), mail.WithPassword(config.Password))
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		Config: config,
		Client: c,
	}

	return mailer, nil
}

// SendMail sends a given mail message.
func (m *Mailer) SendMail(message *MailMessage) error {
	msg := mail.NewMsg()
	if err := msg.From(m.Config.From); err != nil {
		return err
	}
	if err := msg.To(message.To...); err != nil {
		return err
	}
	msg.Subject(message.Subject)
	msg.SetBodyString(mail.ContentType(message.ContentType), message.Body)
	if err := m.Client.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

// SendMailWithContext sends a given mail message with a given context.
func (m *Mailer) SendMailWithContext(message *MailMessage, context context.Context) error {
	msg := mail.NewMsg()
	if err := msg.From(m.Config.From); err != nil {
		return err
	}
	if err := msg.To(message.To...); err != nil {
		return err
	}
	msg.Subject(message.Subject)
	msg.SetBodyString(mail.ContentType(message.ContentType), message.Body)
	if err := m.Client.DialAndSendWithContext(context, msg); err != nil {
		return err
	}
	return nil
}

// Close closes the mailer client.
func (m *Mailer) Close() error {
	return m.Client.Close()
}
