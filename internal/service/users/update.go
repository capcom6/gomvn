package users

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/gomvn/gomvn/internal/entity"
)

func (s *Service) Update(id uint, deploy bool, paths []string) (*entity.User, error) {
	var user entity.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, path := range paths {
			userPath := entity.NewPathID(user.ID, path)
			q := tx.Assign(map[string]any{"Deploy": deploy, "UpdatedAt": now}).
				FirstOrCreate(&userPath)
			if err := q.Error; err != nil {
				return err
			}
		}

		return tx.Model(&user).
			Updates(map[string]any{"UpdatedAt": now}).
			Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
