package service

import (
	"github.com/nozgurozturk/marvin/notifier/entity"
	"github.com/nozgurozturk/marvin/notifier/internal/storage"
)

type RepoService interface {
	FindById(repoId string) (*entity.RepoDTO, error)
}

type repoService struct {
	repository storage.RepoRepository
}

func NewRepoService(r storage.RepoRepository) RepoService {
	return &repoService{
		repository: r,
	}
}

func (s *repoService) FindById(repoId string) (*entity.RepoDTO, error) {
	repo, err := s.repository.FindById(repoId)
	if err != nil {
		return nil, err
	}
	return entity.ToRepoDTO(repo), nil
}
