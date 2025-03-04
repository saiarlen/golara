package settings

import (
	"ekycapp/app/models"
	"ekycapp/config"
	"fmt"
)

type General struct {
	AdminUrl        string
	ClientUrl       string
	CompanyUrl      string
	EntityFullName  string
	EntityShortName string
	Mail            MailCreds
	Sms             SmsCreds
}

type MailCreds struct {
	Host        string
	Port        any
	UserName    string
	Pass        string
	FromName    string
	FromAddress string
	Encryption  string
}

type SmsCreds struct {
	SmsUrl  string
	SmsSid  string
	SmsUser string
	SmsPass string
}
type BackOfficeCreds struct {
	BkUrl         string
	BkTokenGenUrl string
	BkUserName    string
	BkPass        string
	BkGrantType   string
	BkVarOne      string
	BkVarTwo      string
}

// Global variable to cache masters data
var credsCache *General

func GetCreds() (*General, error) {

	if credsCache != nil {
		return credsCache, nil // Return cached data
	}

	credentials := &models.ApxCredentials{}

	err := config.DB.First(credentials).Error
	if err != nil {
		return nil, fmt.Errorf("creds.go: Does not have enough values in creds table %v", err)

	}

	mailConfig, err := getMailConfig(credentials.Email, credentials.ActiveEmail)
	if err != nil {
		return nil, err
	}

	smsConfig, err := getSmsConfig(credentials.Sms, credentials.ActiveSms)
	if err != nil {
		return nil, err
	}

	creds := &General{
		AdminUrl:        credentials.General["ADMIN_URL"].(string),
		ClientUrl:       credentials.General["CLIENT_URL"].(string),
		CompanyUrl:      credentials.General["COMPANY_URL"].(string),
		EntityFullName:  credentials.General["ENTITY_FULL_NAME"].(string),
		EntityShortName: credentials.General["ENTITY_SHORT_NAME"].(string),
		Mail:            *mailConfig,
		Sms:             *smsConfig,
	}
	credsCache = creds
	return credsCache, nil

}

func getMailConfig(email map[string]interface{}, activeEmail string) (*MailCreds, error) {
	activeMail, _ := email[activeEmail].(map[string]interface{})

	mailConfig := &MailCreds{
		Host:        activeMail["MAIL_HOST"].(string),
		Port:        activeMail["MAIL_PORT"],
		UserName:    activeMail["MAIL_USERNAME"].(string),
		Pass:        activeMail["MAIL_PASSWORD"].(string),
		FromName:    activeMail["MAIL_FROM_NAME"].(string),
		FromAddress: activeMail["MAIL_FROM_ADDRESS"].(string),
		Encryption:  activeMail["MAIL_ENCRYPTION"].(string),
	}

	return mailConfig, nil
}

func getSmsConfig(sms map[string]interface{}, activeSms string) (*SmsCreds, error) {
	activeM, _ := sms[activeSms].(map[string]interface{})

	smsConfig := &SmsCreds{
		SmsUrl:  activeM["SMSURL"].(string),
		SmsSid:  activeM["SMSSID"].(string),
		SmsUser: activeM["SMSUSER"].(string),
		SmsPass: activeM["SMSPASSWORD"].(string),
	}

	return smsConfig, nil
}

func ResetCredsCache() {
	mastersCache = nil
}
