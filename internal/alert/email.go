package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds the configuration required to send alert emails.
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

// emailNotifier sends alert notifications via email using SMTP.
type emailNotifier struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailNotifier creates a new Notifier that delivers alerts by email.
// It uses PLAIN auth against the given SMTP server.
func NewEmailNotifier(cfg EmailConfig) Notifier {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	return &emailNotifier{cfg: cfg, auth: auth}
}

// Notify sends the alert as an email to all configured recipients.
func (e *emailNotifier) Notify(a Alert) error {
	if len(e.cfg.To) == 0 {
		return fmt.Errorf("email notifier: no recipients configured")
	}

	subject := fmt.Sprintf("[portwatch] %s on port %d", a.Kind, a.Port)
	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		strings.Join(e.cfg.To, ", "),
		e.cfg.From,
		subject,
		a.String(),
	)

	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	if err := smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, []byte(body)); err != nil {
		return fmt.Errorf("email notifier: send failed: %w", err)
	}
	return nil
}
