package main

import (
	"fmt"

	"github.com/codebasky/lenslocked/model"
)

func main() {
	es := model.NewEmailService(
		model.SMTPConfig{
			Host:     "sandbox.smtp.mailtrap.io",
			Port:     2525,
			User:     "test",
			Password: "test",
		},
	)
	err := es.Send(
		model.Email{
			From:      "support@lenslocked.com",
			To:        "test@gmail.com",
			Subject:   "Test",
			PlainText: "Test mail for email feature",
		},
	)
	if err != nil {
		fmt.Printf("Failed to send email: %s", err)
		return
	}
	fmt.Println("Email sent successfully")
}
