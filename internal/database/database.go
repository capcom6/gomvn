package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gomvn/gomvn/internal/entity"
)

func New() (*gorm.DB, error) {
	if err := os.MkdirAll("data", os.ModeDir); err != nil {
		log.Println("cannot create data directory")
		return nil, err
	}

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// db, err := gorm.Open("sqlite3", "data/data.db")
	db, err := gorm.Open(sqlite.Open("data/data.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	// db.LogMode(true)
	db.AutoMigrate(&entity.User{}, &entity.Path{})
	// ALTER TABLE doesn't work for SQLite
	// db.Model(&entity.Path{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	return db, nil
}
