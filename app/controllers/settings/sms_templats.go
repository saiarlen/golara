package settings

import (
	"ekycapp/app/models"
	"ekycapp/config"
	"fmt"
)

func GetSmsTemp(templateName string) (msg, id string, err error) {

	var smsTemps models.ApxSmsTemplates
	masters, err := GetMasters()
	if err != nil {
		return "", "", err
	}

	err = config.DB.Where("template_name = ? AND service = ?", templateName, masters.Sms).First(&smsTemps).Error
	if err != nil {
		return "", "", fmt.Errorf("%s: sms template record is not available error: %v", templateName, err)
	}

	return smsTemps.TemplateMsg, smsTemps.TemplateId, nil

}
