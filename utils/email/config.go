package email

import (
	"os"
	"strconv"
)

func NewSMTPFromEnv() EmailSender {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return NewSMTPSender(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_FROM"),
	)
}
