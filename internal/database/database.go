package database

import (
	"log"
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
	DRIVER_SQLITE   = "sqlite"
	DRIVER_MYSQL    = "mysql"
	DRIVER_POSTGRES = "postgres"
)

func New(conf *config.App) (*gorm.DB, error) {
	if err := os.MkdirAll("data", os.ModeDir); err != nil {
		log.Println("cannot create data directory")
		return nil, err
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
	case DRIVER_MYSQL:
		dialector = mysql.Open(conf.Database.DSN)
	case DRIVER_POSTGRES:
		dialector = postgres.Open(conf.Database.DSN)
	default:
		dialector = sqlite.Open(conf.Database.DSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: debugLog,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.User{}, &entity.Path{})

	return db, nil
}
