package service

import (
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/storage"
)

type AuthService interface {
	CreateAuth(token *entity.Token) *errors.AppError
	DeleteAuth(uuid string) *errors.AppError
	FindAuth(uuid string) (string, *errors.AppError)
}

type authService struct {
	repository storage.AuthRepository
}

func NewAuthService(r storage.AuthRepository) AuthService {
	return &authService{
		repository: r,
	}
}

func (s *authService) CreateAuth(token *entity.Token) *errors.AppError {

	err := s.repository.CreateAuth(token)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (s *authService) DeleteAuth(uuid string) *errors.AppError {

	err := s.repository.DeleteAuth(uuid)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (r *authService) FindAuth(uuid string) (string, *errors.AppError) {
	id, err := r.repository.FindAuth(uuid)
	if err != nil {
		return "", errors.InternalServer(err.Error())
	}
	return id, nil
}
