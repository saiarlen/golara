package models

import (
	"gorm.io/gorm"
)

type ApxCredentials struct {
	ID          uint                   `gorm:"primaryKey"`
	General     map[string]interface{} `gorm:"column:general;type:json;serializer:json"`
	Email       map[string]interface{} `gorm:"column:email;type:json;serializer:json"`
	Sms         map[string]interface{} `gorm:"column:sms;type:json;serializer:json"`
	ActiveEmail string                 `gorm:"column:active_email"`
	ActiveSms   string                 `gorm:"column:active_sms"`
	gorm.Model
}

type ApxMasters struct {
	ID          uint   `gorm:"primaryKey"`
	OptionKey   string `gorm:"column:option_key;type:varchar(32)"`
	OptionValue string `gorm:"column:option_value;type:longtext"`
	gorm.Model
}
type ApxSmsTemplates struct {
	ID           uint   `gorm:"primaryKey"`
	Service      string `gorm:"column:service;type:varchar(64)"`
	TemplateName string `gorm:"column:template_name;type:varchar(64)"`
	TemplateId   string `gorm:"column:template_id;type:longtext;default:null"`
	TemplateMsg  string `gorm:"column:template_msg;type:longtext;default:null"`
	TemplateType string `gorm:"column:template_type;type:varchar(64);default:null"`
	gorm.Model
}

type ApxEmailTemplates struct {
	ID              uint   `gorm:"primaryKey"`
	TemplateName    string `gorm:"column:template_name;type:varchar(64)"`
	TemplateSubject string `gorm:"column:template_subject;type:longtext"`
	TemplateBody    string `gorm:"column:template_body;type:longtext"`
	TemplateType    string `gorm:"column:template_type;type:varchar(64);default:null"`
	gorm.Model
}
