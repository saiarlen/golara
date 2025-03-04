package global

import (
	"ekycapp/app/controllers/settings"
	"ekycapp/app/providers"

	"fmt"
	"strings"
)

func Notify(mobile, email, templateName string, smsParms, emailParms map[string]string, optvar string) error {

	// Send SMS if mobile is provided
	if mobile != "" {
		// Get SMS template
		msg, id, err := settings.GetSmsTemp(templateName)
		if err != nil {
			return err
		}

		// Replace SMS placeholders with values from smsParms
		if len(smsParms) > 0 {
			for key, value := range smsParms {
				placeholder := fmt.Sprintf("{{$%s}}", key)
				msg = strings.ReplaceAll(msg, placeholder, value)
			}
		}
		fmt.Println(msg)
		_, err = providers.SendSMS(mobile, id, msg, optvar)
		if err != nil {
			return fmt.Errorf("failed to send SMS: %w", err)
		}
	}

	// Send email if email is provided
	if email != "" {
		// Get email template
		emailSub, emailBody, err := settings.GetEmailTemp(templateName)
		if err != nil {
			return fmt.Errorf("failed to get email template: %w", err)
		}

		err = providers.SendEmail(email, emailSub, emailParms, emailBody, optvar)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}
