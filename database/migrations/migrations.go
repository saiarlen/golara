package migrations

import (
	"ekycapp/app/models"
	"ekycapp/config"
)

func Migrations() error {
	db := config.DB

	//php - collation - utf8mb4_unicode_ci
	//go - collation - utf8mb4_0900_ai_ci

	//db.Migrator().AlterColumn(&models.ApxApiUser{}, "kbr")
	//Add All the migrations here

	//Gobal
	db.AutoMigrate(&models.ApxMasters{})

	// Change the order of the "kbr" column to be after the "ckr" column
	//utils.AlterColumnOrder(db, "apx_api_users", "name", "email")

	return nil
}
