package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Status   string `json:"status" gorm:"size:20;default:active"`
}

// TableName returns the table name
func (User) TableName() string {
	return "users"
}