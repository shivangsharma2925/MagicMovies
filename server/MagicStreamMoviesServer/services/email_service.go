package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendEmail(to, subject, html string) error {

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From:    "Magic Movies <noreply@magicmovies.site>",
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	_, err := client.Emails.Send(params)

	return err
}

func SendVerificationOTP(email string, otp string) error {

	body := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Your verification code is:</p>
		<h1>%s</h1>
		<p>This code <b>expires in 10 minutes</b>.</p>
		<p>If this action was not initiated by you, ignore this email.</p>
	`, otp)

	return SendEmail(
		email,
		"Verify your email",
		body,
	)
}

func SendPasswordResetOTP(email string, link string) error {

	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Your reset link is:</p>
		<p>
			<a href="%s">Reset Password</a>
		</p>
		<p>This link <b>expires in 15 minutes</b>.</p>
		<p>If this action was not initiated by you, ignore this email.</p>
	`, link)

	return SendEmail(
		email,
		"Password Reset",
		body,
	)
}
