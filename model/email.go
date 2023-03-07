package model

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}
type EmailService struct {
	dialer *gomail.Dialer
}

func NewEmailService(cfg SMTPConfig) *EmailService {
	return &EmailService{
		dialer: gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password),
	}
}

func (es *EmailService) Send(email Email) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", email.From)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	switch {
	case email.PlainText != "" && email.HTML != "":
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	case email.PlainText != "":
		msg.SetBody("text/plain", email.PlainText)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}
	return es.dialer.DialAndSend(msg)
}

func (es *EmailService) ForgotPassword(to string, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		From:      DefaultSender,
		To:        to,
		PlainText: "To reset your password, please visit the following link: " + resetURL,
		HTML:      `<p>To reset your password, please visit the following link: <a href="` + resetURL + `">` + resetURL + `</a></p>`,
	}
	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email send error: %w", err)
	}
	return nil
}
