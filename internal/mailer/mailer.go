package mailer

import (
	"github.com/wneessen/go-mail"
)

type MailConfig struct {
	From     string
	Username string
	Password string
	Port     int
	Host     string
}

type Mailer struct {
	*mail.Client
	config *MailConfig
}

// ConnectToMailer will create a new mail client and return it.
func ConnectToMailer(config *MailConfig) (*Mailer, error) {
	c, err := mail.NewClient(config.Host, mail.WithPort(config.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.Username), mail.WithPassword(config.Password))
	if err != nil {
		return nil, err
	}
	return &Mailer{c, config}, nil
}

// SendMail will send a mail using the given mail client.
func (m *Mailer) SendMail(to string, subject string, body string) error {
	msg := mail.NewMsg()
	if err := msg.From(m.config.From); err != nil {
		return err
	}
	if err := msg.To(to); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)
	if err := m.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
