package database

import (
	"fmt"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Migration string    `gorm:"size:255;uniqueIndex"`
	Batch     int       `gorm:"index"`
	CreatedAt time.Time
}

// MigrationFile represents a migration file
type MigrationFile struct {
	Name string
	Up   func(*gorm.DB) error
	Down func(*gorm.DB) error
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationFile
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]MigrationFile, 0),
	}
}

// Add adds a migration to the migrator
func (m *Migrator) Add(name string, up, down func(*gorm.DB) error) {
	m.migrations = append(m.migrations, MigrationFile{
		Name: name,
		Up:   up,
		Down: down,
	})
}

// Run executes pending migrations
func (m *Migrator) Run() error {
	// Create migrations table if not exists
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current batch number
	var lastMigration Migration
	m.db.Order("batch DESC").First(&lastMigration)
	nextBatch := lastMigration.Batch + 1

	// Get executed migrations
	var executedMigrations []Migration
	m.db.Find(&executedMigrations)
	
	executedMap := make(map[string]bool)
	for _, migration := range executedMigrations {
		executedMap[migration.Migration] = true
	}

	// Sort migrations by name (timestamp)
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Name < m.migrations[j].Name
	})

	// Run pending migrations
	for _, migration := range m.migrations {
		if !executedMap[migration.Name] {
			log.Printf("Running migration: %s", migration.Name)
			
			if err := migration.Up(m.db); err != nil {
				return fmt.Errorf("migration %s failed: %w", migration.Name, err)
			}

			// Record migration
			migrationRecord := Migration{
				Migration: migration.Name,
				Batch:     nextBatch,
				CreatedAt: time.Now(),
			}
			
			if err := m.db.Create(&migrationRecord).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
			}
			
			log.Printf("✅ Migration completed: %s", migration.Name)
		}
	}

	return nil
}

// Rollback rolls back the last batch of migrations
func (m *Migrator) Rollback(steps int) error {
	if steps <= 0 {
		steps = 1
	}

	// Get migrations to rollback
	var migrations []Migration
	m.db.Order("batch DESC, id DESC").Limit(steps).Find(&migrations)

	if len(migrations) == 0 {
		log.Println("No migrations to rollback")
		return nil
	}

	// Group by batch and rollback
	batchMap := make(map[int][]Migration)
	for _, migration := range migrations {
		batchMap[migration.Batch] = append(batchMap[migration.Batch], migration)
	}

	// Get batch numbers and sort descending
	var batches []int
	for batch := range batchMap {
		batches = append(batches, batch)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(batches)))

	for _, batch := range batches {
		batchMigrations := batchMap[batch]
		
		// Reverse order for rollback
		for i := len(batchMigrations) - 1; i >= 0; i-- {
			migration := batchMigrations[i]
			
			// Find migration file
			var migrationFile *MigrationFile
			for _, mf := range m.migrations {
				if mf.Name == migration.Migration {
					migrationFile = &mf
					break
				}
			}

			if migrationFile == nil {
				log.Printf("⚠️  Migration file not found: %s", migration.Migration)
				continue
			}

			log.Printf("Rolling back migration: %s", migration.Migration)
			
			if err := migrationFile.Down(m.db); err != nil {
				return fmt.Errorf("rollback %s failed: %w", migration.Migration, err)
			}

			// Remove migration record
			if err := m.db.Delete(&migration).Error; err != nil {
				return fmt.Errorf("failed to remove migration record %s: %w", migration.Migration, err)
			}
			
			log.Printf("✅ Rollback completed: %s", migration.Migration)
		}
	}

	return nil
}

// Status shows migration status
func (m *Migrator) Status() error {
	var executedMigrations []Migration
	m.db.Order("batch ASC, id ASC").Find(&executedMigrations)
	
	executedMap := make(map[string]Migration)
	for _, migration := range executedMigrations {
		executedMap[migration.Migration] = migration
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("================")
	
	for _, migration := range m.migrations {
		if executed, exists := executedMap[migration.Name]; exists {
			fmt.Printf("✅ %s (Batch: %d)\n", migration.Name, executed.Batch)
		} else {
			fmt.Printf("❌ %s (Pending)\n", migration.Name)
		}
	}
	
	return nil
}