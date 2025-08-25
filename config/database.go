package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {
	DATABASE_URI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		Denv("DB_USERNAME"), Denv("DB_PASSWORD"), Denv("DB_HOST"), Denv("DB_PORT"), Denv("DB_DATABASE"))

	// Create logs directory if not exists
	if err := os.MkdirAll("storage/logs", 0755); err != nil {
		log.Printf("Warning: Could not create logs directory: %v", err)
	}

	// Setup logger
	logFile, err := os.OpenFile("storage/logs/gorm.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Warning: Could not open log file: %v", err)
		logFile = os.Stdout
	}

	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	var loggerConfig logger.Interface
	if Denv("APP_ENV") == "production" {
		loggerConfig = newLogger
	} else {
		loggerConfig = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 loggerConfig,
	})

	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// Test connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database Connected Successfully")
	return nil
}