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
const otpLength = 8

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := make([]byte, otpLength)

	// Ensure at least 3 digits
	digitCount := 0
	for digitCount < 3 {
		pos := r.Intn(otpLength)
		if otp[pos] == 0 { // Ensure we don't overwrite already set digits
			otp[pos] = otpChars[r.Intn(10)+26] // Digits are in the last 10 characters of otpChars
			digitCount++
		}
	}

	// Fill the remaining positions
	for i := range otp {
		if otp[i] == 0 {
			otp[i] = otpChars[r.Intn(len(otpChars))]
		}
	}

	return string(otp)
}

func SendOTPEmail(to, otp string) error {
	// Prepare the email data
	data := template.OTPEmailData{
		OTP: otp,
	}

	// Generate the email body
	body, err := template.GenerateOTPEmail(data)
	if err != nil {
		log.Printf("Error generating OTP email: %v", err)
		return err
	}

	// Set up email subject and recipient
	subject := "Verification Code"
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
