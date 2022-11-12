package service

import (
	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/entity"
	"github.com/gomvn/gomvn/internal/service/storage"
)

func NewRepoService(conf *config.App, storage *storage.LocalStorage, ps *PathService) *RepoService {
	return &RepoService{
		repository: conf.Repository,
		storage:    storage,
		ps:         ps,
	}
}

type RepoService struct {
	repository []string
	storage    *storage.LocalStorage
	ps         *PathService
}

func (s *RepoService) GetRepositories() map[string][]*entity.Artifact {
	result := map[string][]*entity.Artifact{}
	for _, repo := range s.repository {
		result[repo] = s.storage.List(repo)
	}
	return result
}
