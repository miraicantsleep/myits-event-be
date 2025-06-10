package utils

import (
	"bytes" // Required for parsing HTML template
	"html/template" // Required for parsing HTML template
	"io"            // Required for gomail.SetCopyFunc

	// "path/filepath" // If loading template from specific path structure beyond just filename
	"github.com/miraicantsleep/myits-event-be/config"
	"gopkg.in/gomail.v2"
)

func SendMail(toEmail string, subject string, body string) error {
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.SenderName)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthUsername,
		emailConfig.AuthPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

// SendInvitationMail sends a styled HTML invitation email with an embedded QR code.
func SendInvitationMail(toEmail string, subject string, templateData map[string]interface{}, qrCodeImage []byte) error {
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		return err
	}

	// Parse the HTML email template
	// Ensure the path to template is correct. Assuming it's relative to where the binary runs,
	// or adjust path as needed e.g. using an absolute path or relative to GOPATH/module root.
	// For simplicity, let's assume it's in a known relative path for now.
	// This path might need to be configurable or determined more robustly in a real app.
	tmpl, err := template.ParseFiles("utils/email-template/invitation_mail.html")
	if err != nil {
		return err // Could not parse template
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		return err // Could not execute template
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.SenderName)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)

	// Embed the QR code image
	if len(qrCodeImage) > 0 { // qrCodeImage is []byte
		// The first argument to Embed, "qr_code_image", will be used as the default for Content-ID
		// if not overridden by SetHeader. It's also used as the "filename" in Content-Disposition.
		// We explicitly set Content-ID to ensure it's exactly "<qr_code_image>" to match "cid:qr_code_image" in HTML.
		mailer.Embed("qr_code_image.png", // This name is primarily for reference/filename hint
			gomail.SetHeader(map[string][]string{
				"Content-Type": {"image/png"},
				"Content-ID":   {"<qr_code_image>"}, // Crucial for <img src="cid:qr_code_image">
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, errWrite := w.Write(qrCodeImage)
				return errWrite
			}),
		)
	}

	mailer.SetBody("text/html", body.String())

	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthUsername,
		emailConfig.AuthPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
