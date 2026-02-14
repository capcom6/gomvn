package database

import (
	"fmt"
	"log" //nolint:depguard // TODO
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/entity"
)

const (
	DriverSQLite   = "sqlite"
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
)

func New(conf *config.App) (*gorm.DB, error) {
	if err := os.MkdirAll("data", os.ModeDir); err != nil {
		log.Println("cannot create data directory")
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	var debugLog logger.Interface
	if conf.Debug {
		debugLog = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
	}

	var dialector gorm.Dialector
	switch conf.Database.Driver {
	case DriverMySQL:
		dialector = mysql.Open(conf.Database.DSN)
	case DriverPostgres:
		dialector = postgres.Open(conf.Database.DSN)
	case DriverSQLite:
		dialector = sqlite.Open(conf.Database.DSN)
	default:
		dialector = sqlite.Open("data/data.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: debugLog,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if migrateErr := db.AutoMigrate(new(entity.User), new(entity.Path)); migrateErr != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", migrateErr)
	}

	return db, nil
}
