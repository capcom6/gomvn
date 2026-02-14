package users

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gomvn/gomvn/internal/entity"
	gonanoid "github.com/matoous/go-nanoid"
)

func (s *Service) Create(name string, isAdmin bool, deploy bool, paths []string) (*entity.User, string, error) {
	token, err := gonanoid.ID(TokenLength)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), BcryptCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash token: %w", err)
	}

	now := time.Now()
	user := entity.User{
		Name:      name,
		Admin:     isAdmin,
		TokenHash: string(tokenHash),
		ID:        0,
		CreatedAt: now,
		UpdatedAt: now,
		Paths:     []entity.Path{},
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if crErr := tx.Create(&user).Error; crErr != nil {
			return crErr
		}
		for _, path := range paths {
			userPath := entity.NewPath(user.ID, path, deploy)
			if pathErr := tx.Create(&userPath).Error; pathErr != nil {
				return pathErr
			}
			user.Paths = append(user.Paths, userPath)
		}
		return nil
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	return &user, token, nil
}
