package migrations

import (
	"github.com/test/myapp/framework/database"

	"gorm.io/gorm"
)

// RegisterMigrations registers all migrations
func RegisterMigrations(migrator *database.Migrator) {
	// Register all migrations here
	name, up, down := CreateUsersTableMigration()
	migrator.Add(name, up, down)

	// Add more migrations here as you create them
	// name2, up2, down2 := CreateProductsTableMigration()
	// migrator.Add(name2, up2, down2)
}

// RunMigrations runs all registered migrations
func RunMigrations(db *gorm.DB) error {
	migrator := database.NewMigrator(db)
	RegisterMigrations(migrator)
	return migrator.Run()
}

// RollbackMigrations rolls back migrations
func RollbackMigrations(db *gorm.DB, steps int) error {
	migrator := database.NewMigrator(db)
	RegisterMigrations(migrator)
	return migrator.Rollback(steps)
}

// MigrationStatus shows migration status
func MigrationStatus(db *gorm.DB) error {
	migrator := database.NewMigrator(db)
	RegisterMigrations(migrator)
	return migrator.Status()
}
