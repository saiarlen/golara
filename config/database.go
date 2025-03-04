package config

import (
	"ekycapp/utils"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var errDB error

func ConnectDB() error {
	DATABASE_URI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		Denv("DB_USERNAME"), Denv("DB_PASSWORD"), Denv("DB_HOST"), Denv("DB_PORT"), Denv("DB_DATABASE"))

	newLogger := logger.New(
		log.New(utils.LogFile("gorm.log"), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
		},
	)

	var loggerConfig logger.Interface
	if Env("APP_ENV") == "production" {
		loggerConfig = newLogger
	} else {
		loggerConfig = logger.Default.LogMode(logger.Info)
	}

	DB, errDB = gorm.Open(mysql.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 loggerConfig,
	})

	if errDB != nil {
		fmt.Println("database connection failed: Config error")
		return fmt.Errorf("database connection failed")
	}

	fmt.Println("Database Connected Successfully")

	return nil
}
