package adapter

import (
	"fmt"
	"net/smtp"
)

type smtpEmailSender struct {
	from     string
	password string
	host     string
	port     int
}

func New(from, password, host string, port int) *smtpEmailSender {
	return &smtpEmailSender{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (s *smtpEmailSender) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	err := smtp.SendMail(addr, auth, s.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send SMTP email: %w", err)
	}

	return nil
}
