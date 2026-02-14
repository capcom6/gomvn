package users

import (
	"fmt"
	"log" //nolint:depguard // TODO

	"gorm.io/gorm"

	"github.com/gomvn/gomvn/internal/entity"
)

func Initialize(db *gorm.DB, us *Service) error {
	var count int64
	if err := db.Model((*entity.User)(nil)).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count == 0 {
		log.Println("Initializing first user...")

		user, token, err := us.Create("admin", true, true, []string{"/"})
		if err != nil {
			return err
		}

		log.Printf("USER ID: %d\n", user.ID)
		log.Printf("USERNAME: %s\n", user.Name)
		log.Printf("TOKEN: %s\n", token)
	}

	return nil
}
