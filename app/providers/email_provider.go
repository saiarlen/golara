package providers

import (
	"ekycapp/app/controllers/settings"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, data map[string]string, body string, attachmentPath string) error {
	creds, err := settings.GetCreds()
	if err != nil {
		return err
	}

	if to == "" {
		return fmt.Errorf("empty email address: func:SendEmail()")
	}

	smtpHost := creds.Mail.Host
	smtpPort := cast.ToInt(creds.Mail.Port)
	senderEmail := creds.Mail.UserName
	senderPassword := creds.Mail.Pass
	fromName := creds.Mail.FromName
	fromAddress := creds.Mail.FromAddress

	// If data is provided, parse placeholders in the body template
	if data != nil {
		body = ParseTemplate(body, data)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", message.FormatAddress(fromAddress, fromName))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Attach file if an attachment path is provided
	if attachmentPath != "" {
		message.Attach(attachmentPath)
	}

	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)
	return dialer.DialAndSend(message)
}

func ParseTemplate(template string, data map[string]string) string {
	for key, value := range data {
		placeholder := "{{$" + key + "}}"
		template = strings.Replace(template, placeholder, value, -1)
	}
	return template
}
