package users

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gomvn/gomvn/internal/entity"
)

func (s *Service) UpdateToken(id uint) (*entity.User, string, error) {
	var user entity.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	token, err := gonanoid.ID(TokenLength)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), BcryptCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash token: %w", err)
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).
			Updates(map[string]any{
				"TokenHash": string(tokenHash),
				"UpdatedAt": time.Now(),
			}).
			Error
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to update token: %w", err)
	}

	return &user, token, nil
}
