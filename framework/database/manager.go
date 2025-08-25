package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager manages multiple database connections
type DatabaseManager struct {
	connections map[string]*gorm.DB
	default_    string
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		connections: make(map[string]*gorm.DB),
	}
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	SSLMode  string `yaml:"sslmode"`
	Timezone string `yaml:"timezone"`
}

// Connect establishes database connection
func (dm *DatabaseManager) Connect(name string, config DatabaseConfig) error {
	var dialector gorm.Dialector
	
	switch config.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, 
			config.Database, config.Charset)
		dialector = mysql.Open(dsn)
		
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			config.Host, config.Username, config.Password, config.Database, 
			config.Port, config.SSLMode, config.Timezone)
		dialector = postgres.Open(dsn)
		
	case "sqlite":
		dialector = sqlite.Open(config.Database)
		
	default:
		return fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to %s database: %w", config.Driver, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dm.connections[name] = db
	
	if dm.default_ == "" {
		dm.default_ = name
	}
	
	log.Printf("✅ Connected to %s database: %s", config.Driver, name)
	return nil
}

// Connection returns a database connection by name
func (dm *DatabaseManager) Connection(name ...string) *gorm.DB {
	connName := dm.default_
	if len(name) > 0 {
		connName = name[0]
	}
	
	if conn, exists := dm.connections[connName]; exists {
		return conn
	}
	
	log.Printf("⚠️  Database connection '%s' not found, using default", connName)
	return dm.connections[dm.default_]
}

// SetDefault sets the default connection
func (dm *DatabaseManager) SetDefault(name string) {
	dm.default_ = name
}

// Close closes all database connections
func (dm *DatabaseManager) Close() error {
	for name, db := range dm.connections {
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection %s: %v", name, err)
		}
	}
	return nil
}