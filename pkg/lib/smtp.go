package lib

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var (
	auth   smtp.Auth
	sender string
)

func InitSMTP() {
	// Get authentication information from environment variables
	sender = os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if sender == "" || password == "" {
		log.Fatal("SMTP authentication information not set")
	}

	// Set up authentication information for Gmail SMTP
	auth = smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	if auth == nil {
		log.Fatal("SMTP authentication information is invalid")
	} else {
		log.Println("SMTP authentication information set")
	}
}

func SendEmail(to []string, subject string, body string) error {
	// Compose the message with headers for HTML content
	msg := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", to[0], subject, body)

	// Send the email
	err := smtp.SendMail("smtp.gmail.com:587", auth, sender, to, []byte(msg))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %s", to[0])
	return nil
}
