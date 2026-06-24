package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"laundry-system/internal/config"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendEmail(to []string, subject, htmlBody string) error {
	log.Printf("[EMAIL] Sending email to %v with subject: %s", to, subject)

	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(to, ", "),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
		"Date":         time.Now().Format(time.RFC1123Z),
	}

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	auth := s.getAuth()
	addr := s.config.Host + ":" + s.config.Port

	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, s.config.FromEmail, to, []byte(message.String()))
	}
	return smtp.SendMail(addr, auth, s.config.FromEmail, to, []byte(message.String()))
}

func (s *EmailService) getAuth() smtp.Auth {
	if s.config.AuthMethod == "LOGIN" {
		return &loginAuth{username: s.config.Username, password: s.config.Password}
	}
	return smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
}

func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{ServerName: s.config.Host}); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT command failed for %s: %w", recipient, err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}
	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("[EMAIL] Email sent successfully to %v", to)
	return client.Quit()
}

// loginAuth implements LOGIN authentication for legacy SMTP servers.
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(strings.TrimSpace(string(fromServer))) {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server prompt: %s", fromServer)
		}
	}
	return nil, nil
}
