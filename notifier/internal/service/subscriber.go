package service

import (
	"github.com/nozgurozturk/marvin/notifier/entity"
	"github.com/nozgurozturk/marvin/notifier/internal/storage"
)

type SubscriberService interface {
	FindByEmailAndRepoId(email string, repoId string) (*entity.SubscriberDTO, error)
	FindById(subId string) (*entity.SubscriberDTO, error)
	FindAll(repoId string) ([]*entity.SubscriberDTO, error)
	GetAll() ([]*entity.SubscriberDTO, error)
}

type subService struct {
	repository storage.SubscriberRepository
}

func NewSubscriberService(r storage.SubscriberRepository) SubscriberService {
	return &subService{
		repository: r,
	}
}

func (s *subService) FindByEmailAndRepoId(email string, repoId string) (*entity.SubscriberDTO, error) {
	subscriber, err := s.repository.FindByEmailAndRepoId(email, repoId)
	if err != nil {
		return nil, err
	}
	return entity.ToSubscriberDTO(subscriber), nil
}

func (s *subService) FindById(subId string) (*entity.SubscriberDTO, error) {
	subscriber, err := s.repository.FindById(subId)
	if err != nil {
		return nil, err
	}
	return entity.ToSubscriberDTO(subscriber), nil
}

func (s *subService) GetAll() ([]*entity.SubscriberDTO, error) {
	subs, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return entity.ToSubscriberDTOs(subs), nil
}

func (s *subService) FindAll(repoId string) ([]*entity.SubscriberDTO, error) {
	subs, err := s.repository.FindAll(repoId)
	if err != nil {
		return nil, err
	}

	return entity.ToSubscriberDTOs(subs), nil
}
