package user

import (
	"gorm.io/gorm"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/entity"
)

const (
	TokenLength = 36
	BcryptCost  = 12
)

func New(db *gorm.DB, conf *config.App) *Service {
	return &Service{
		db:   db,
		conf: conf,
	}
}

type Service struct {
	db   *gorm.DB
	conf *config.App
}

func (s *Service) GetAll(limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	if err := s.db.Limit(limit).Offset(offset).Preload("Paths").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	if err := s.db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (s *Service) GetByName(name string) (*entity.User, error) {
	var user entity.User
	if err := s.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetPaths(user *entity.User) ([]entity.Path, error) {
	var paths []entity.Path
	q := s.db.Model(&entity.Path{}).Where("user_id = ?", user.ID).Find(&paths)
	if err := q.Error; err != nil {
		return nil, err
	}
	return paths, nil
}

func (s *Service) GetDefaultPermissions() *config.Permissions {
	return &s.conf.Permissions
}
