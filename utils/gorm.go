package utils

import (
	"fmt"

	"gorm.io/gorm"
)

// AlterColumnOrder modifies the order of a column in a table (separate function)
func AlterColumnOrder(db *gorm.DB, tableName, firstColumn, afterColumn string) error {

	// Construct and execute ALTER TABLE statement

	getSchemaType := fmt.Sprintf("SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE table_name = '%s' AND column_name = '%s'",
		tableName, firstColumn)

	var dataType string
	row := db.Raw(getSchemaType).Row()
	err := row.Scan(&dataType)

	if err != nil {
		// Handle error if data type cannot be scanned
		return err
	}

	definition := dataType

	sql := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s AFTER %s",
		tableName, firstColumn, definition, afterColumn)

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to alter column order: %w", err)
	}

	fmt.Println("Successfully altered column order of", firstColumn, "in table", tableName)
	return nil
}
