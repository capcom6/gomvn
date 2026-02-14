package service

import (
	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/entity"
	"github.com/gomvn/gomvn/internal/service/storage"
)

type RepoService struct {
	repository []string
	storage    *storage.Storage
	ps         *PathService
}

func NewRepoService(conf *config.App, storage *storage.Storage, ps *PathService) *RepoService {
	return &RepoService{
		repository: conf.Repository,
		storage:    storage,
		ps:         ps,
	}
}

func (s *RepoService) GetRepositories() map[string][]*entity.Artifact {
	result := map[string][]*entity.Artifact{}
	for _, repo := range s.repository {
		result[repo], _ = s.storage.ListArtifacts(repo)
	}
	return result
}
