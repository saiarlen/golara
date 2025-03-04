package settings

import (
	"ekycapp/app/models"
	"ekycapp/config"
	"fmt"
)

func GetEmailTemp(templateName string) (subject string, body string, err error) {

	var emailTemps models.ApxEmailTemplates

	err = config.DB.Where("template_name = ?", templateName).First(&emailTemps).Error
	if err != nil {
		return "", "", fmt.Errorf(templateName+": email template record is not available error:%v", err)

	}

	return emailTemps.TemplateSubject, emailTemps.TemplateBody, nil
}
