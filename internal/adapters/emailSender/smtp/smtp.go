package adapter

import (
	"fmt"
	"net/smtp"
)

// smtpEmailSender is a struct that holds configuration for sending emails via SMTP.
type smtpEmailSender struct {
	from     string // Email address used as the sender
	password string // SMTP password or app-specific password for authentication
	host     string // SMTP server host (e.g., smtp.gmail.com)
	port     int    // SMTP server port (e.g., 587 for TLS)
}

// New creates a new instance of smtpEmailSender with the given SMTP configuration.
func New(from, password, host string, port int) *smtpEmailSender {
	return &smtpEmailSender{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

// SendEmail sends an HTML email to the specified recipient using SMTP.
func (s *smtpEmailSender) SendEmail(to, subject, body string) error {
	// Set up authentication using the sender's email and password.
	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	// Construct the MIME-compliant email message.
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	// Build the full address (host:port) of the SMTP server.
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Attempt to send the email using smtp.SendMail.
	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send SMTP email: %w", err)
	}

	// Return nil if the email was sent successfully.
	return nil
}
