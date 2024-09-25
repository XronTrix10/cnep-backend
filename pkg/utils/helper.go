package utils

import (
	"math/rand"
	"time"

	"fmt"
	"log"

	"cnep-backend/pkg/lib"
	"cnep-backend/pkg/template"
)

const otpChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := make([]byte, 8)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}
	return string(otp)
}

func SendOTPEmail(to, name, otp string) error {
	// Prepare the email data
	data := template.OTPEmailData{
		Name: name,
		OTP:  otp,
	}

	// Generate the email body
	body, err := template.GenerateOTPEmail(data)
	if err != nil {
		log.Printf("Error generating OTP email: %v", err)
		return err
	}

	// Set up email subject and recipient
	subject := "Your OTP for Authentication"
	recipient := []string{to}

	// Send the email
	err = lib.SendEmail(recipient, subject, body)
	if err != nil {
		log.Printf("Error sending OTP email: %v", err)
		return err
	}

	fmt.Printf("OTP email sent successfully to %s\n", to)
	return nil
}
