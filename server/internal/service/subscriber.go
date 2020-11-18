package service

import (
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/storage"
	"time"
)

type SubscriberService interface {
	// Create creates new subscriber and saves into store
	Create(email string, repoID string) (*entity.SubscriberDTO, *errors.AppError)
	// FindByID returns subscriber with matching id
	FindByID(subID string) (*entity.SubscriberDTO, *errors.AppError)
	// FindByEmailAndRepoID returns subscriber with matching email and repository id
	FindByEmailAndRepoID(email string, repoID string) (*entity.SubscriberDTO, *errors.AppError)
	// FindAll returns subscriber belongs to git repository
	FindAll(repoID string) ([]*entity.SubscriberDTO, *errors.AppError)
	// GetAll returns all subscribers from store
	All() ([]*entity.SubscriberDTO, *errors.AppError)
	// Update insert updated subscriber values into store
	Update(subscriberDTO *entity.SubscriberDTO) (*entity.SubscriberDTO, *errors.AppError)
	// Confirm verify subscriber's email or Unsubscribe
	Confirm(subDTO *entity.SubscriberDTO, confirmed bool) *errors.AppError
	// Delete remove subscriber entity from store
	Delete(subID string) *errors.AppError
	// Delete all subscribers belongs to git repository
	DeleteAll(repoID string) *errors.AppError
}

type subService struct {
	repository storage.SubscriberRepository
}

func NewSubscriberService(r storage.SubscriberRepository) SubscriberService {
	return &subService{
		repository: r,
	}
}

func (s *subService) Create(email string, repoID string) (*entity.SubscriberDTO, *errors.AppError) {
	// creates subscriber dto with default value
	// default notify value
	/*
		{
			frequency: "day",
			weekday: current week day as an integer,
			hour: current hour, minute: current minute
		}
	*/
	now := time.Now()
	subscriber := &entity.SubscriberDTO{
		Email:  email,
		RepoID: repoID,
		Notify: &entity.Notify{
			Frequency: entity.Day,
			Weekday:   now.Weekday(),
			Hour:      now.Hour(),
			Minute:    now.Minute(),
		},
	}

	createdSub, err := s.repository.Create(entity.ToSubscriber(subscriber))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToSubscriberDTO(createdSub), nil
}

func (s *subService) Update(subscriberDTO *entity.SubscriberDTO) (*entity.SubscriberDTO, *errors.AppError) {

	updated, err := s.repository.Update(entity.ToSubscriber(subscriberDTO))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToSubscriberDTO(updated), nil
}

func (s *subService) Confirm(subDTO *entity.SubscriberDTO, confirmed bool) *errors.AppError {

	err := s.repository.Confirm(*subDTO.ID, confirmed)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (s *subService) FindByEmailAndRepoID(email string, repoID string) (*entity.SubscriberDTO, *errors.AppError) {

	subscriber, err := s.repository.FindByEmailAndRepoID(email, repoID)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	if subscriber == nil {
		return nil, errors.NotFound("Subscriber is not found")
	}

	return entity.ToSubscriberDTO(subscriber), nil
}

func (s *subService) FindByID(subID string) (*entity.SubscriberDTO, *errors.AppError) {

	subscriber, err := s.repository.FindByID(subID)
	if err != nil {
		return nil, errors.NotFound("Subscriber is not found")
	}

	if subscriber == nil {
		return nil, errors.NotFound("Subscriber is not found")
	}

	return entity.ToSubscriberDTO(subscriber), nil
}

func (s *subService) All() ([]*entity.SubscriberDTO, *errors.AppError) {

	subs, err := s.repository.GetAll()
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToSubscriberDTOs(subs), nil
}

func (s *subService) FindAll(repoID string) ([]*entity.SubscriberDTO, *errors.AppError) {

	subs, err := s.repository.FindAll(repoID)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToSubscriberDTOs(subs), nil
}

func (s *subService) Delete(subID string) *errors.AppError {

	err := s.repository.Delete(subID)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (s *subService) DeleteAll(repoID string) *errors.AppError {

	err := s.repository.DeleteMany(repoID)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}
