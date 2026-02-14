package users

import (
	"fmt"

	"github.com/gomvn/gomvn/internal/entity"
)

func (s *Service) Delete(id uint) error {
	err := s.db.Where("id = ?", id).Delete(new(entity.User)).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
