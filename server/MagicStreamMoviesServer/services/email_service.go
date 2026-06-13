package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationOTP(email string, otp string) error {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_HOST"),
	)

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	subject := "Subject: Verify your email\n"

	body := fmt.Sprintf("<h2>Email Verification</h2>\n<p>Your verification code is:</p>\n<h1>%s</h1>\n<p>This code <b>expires in 10 minutes</b>.</p>\n<p>If this action is not taken by you, Ignore this email.</p>", otp)

	message := []byte(subject + mime + body)

	addr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))

	return smtp.SendMail(addr, auth, os.Getenv("SMTP_USER"), []string{email}, message)

}

func SendPasswordResetOTP(email string, link string) error {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_HOST"),
	)

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	subject := "Subject: Password Reset\n"

	body := fmt.Sprintf("<h2>Password Reset Request</h2>\n<p>Your reset link is:</p>\n<h4>%s</h4>\n<p>This link <b>expires in 15 minutes</b>.</p>\n<p>If this action is not taken by you, Ignore this email.</p>", link)

	message := []byte(subject + mime + body)

	addr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))

	return smtp.SendMail(addr, auth, os.Getenv("SMTP_USER"), []string{email}, message)
}
