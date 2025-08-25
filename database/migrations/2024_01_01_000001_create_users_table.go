package migrations

import (
	"gorm.io/gorm"
)

// CreateUsersTableMigration creates users table
func CreateUsersTableMigration() (string, func(*gorm.DB) error, func(*gorm.DB) error) {
	return "2024_01_01_000001_create_users_table",
		// Up
		func(db *gorm.DB) error {
			return db.Exec(`
				CREATE TABLE users (
					id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					name VARCHAR(100) NOT NULL,
					email VARCHAR(100) NOT NULL UNIQUE,
					password VARCHAR(255) NOT NULL,
					status VARCHAR(20) DEFAULT 'active',
					created_at TIMESTAMP NULL,
					updated_at TIMESTAMP NULL,
					deleted_at TIMESTAMP NULL,
					INDEX idx_users_deleted_at (deleted_at),
					INDEX idx_users_email (email),
					INDEX idx_users_status (status)
				)
			`).Error
		},
		// Down
		func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS users").Error
		}
}