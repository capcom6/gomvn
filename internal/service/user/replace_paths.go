// Copyright 2022 Aleksandr Soloshenko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"time"

	"github.com/gomvn/gomvn/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) ReplacePaths(id uint, paths []entity.Path) ([]entity.Path, error) {
	var user entity.User
	if err := s.db.Model(&user).Preload("Paths").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	deletePaths := make(map[string]bool)
	for _, path := range user.Paths {
		deletePaths[path.Path] = true
	}

	user.UpdatedAt = time.Now()
	user.Paths = make([]entity.Path, len(paths))
	for i, path := range paths {
		delete(deletePaths, path.Path)

		path.UserID = user.ID
		path.CreatedAt = time.Now()
		path.UpdatedAt = time.Now()
		user.Paths[i] = path
	}

	toDelete := make([]entity.Path, 0)
	for path := range deletePaths {
		toDelete = append(toDelete, entity.Path{
			UserID: id,
			Path:   path,
		})
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if len(toDelete) > 0 {
			if err := tx.Delete(&toDelete).Error; err != nil {
				return err
			}
		}

		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"deploy"}),
		}).Save(&user.Paths).Error; err != nil {
			return err
		}

		return tx.Select("UpdatedAt").Save(&user).Error
	},
	)
	if err != nil {
		return nil, err
	}

	return user.Paths, nil
}
