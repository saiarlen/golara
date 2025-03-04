package models

import (
	"gorm.io/gorm"
)

type ApxApiUser struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Email    string
	Kbr      bool
	Name     string
	Password string
	Ckr      string
	gorm.Model
}
