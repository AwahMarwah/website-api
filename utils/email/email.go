package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type smtpEmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSMTPSender(host string, port int, username, password, from string) EmailSender {
	return &smtpEmailSender{
		Host: host, Port: port, Username: username, Password: password, From: from,
	}
}

func (s *smtpEmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
