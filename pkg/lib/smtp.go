package lib

import (
	"log"
	"net/smtp"
	"os"
)

var auth smtp.Auth

func InitSMTP() {
	// Set up authentication information.
	sender := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if sender == "" || password == "" {
		panic("SMTP authentication information not set")
	}

	// Set up authentication information.
	auth = smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	if auth == nil {
		panic("SMTP authentication information is invalid")
	} else {
		log.Println("SMTP authentication information set")
	}
}

func SendEmail(from string, to []string, subject string, body string) error {

	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail("mail.example.com:25", auth, "sender@example.org", to, msg)
	if err != nil {
		log.Fatal("Error sending email: " + err.Error())
		return err
	}
	return nil
}
