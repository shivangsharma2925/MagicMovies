package services

import (
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailService struct {
	client *resend.Client
	from   string
}

func NewEmailService() *EmailService {

	return &EmailService{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
		from:   os.Getenv("EMAIL_FROM"),
	}
}

func (e *EmailService) SendVerificationOTP(email string, otp string,) error {

	params := &resend.SendEmailRequest{
		From:    e.from,
		To:      []string{email},
		Subject: "Verify your email",
		Html: `
			<h2>Email Verification</h2>
			<p>Your verification code is:</p>
			<h1>` + otp + `</h1>
			<p>This code expires in 10 minutes.</p>
		`,
	}

	_, err := e.client.Emails.Send(params)

	return err
}

func (e *EmailService) SendPasswordResetOTP(email string, link string,) error {

	params := &resend.SendEmailRequest{
		From:    e.from,
		To:      []string{email},
		Subject: "Password Reset",
		Html: `
			<h2>Password Reset Request</h2>
			<p>Your reset link is:</p>
			<h4>` + link + `</h4>
			<p>This link expires in 15 minutes.</p>
		`,
	}

	_, err := e.client.Emails.Send(params)

	return err
}
