package service

import (
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/storage"
)

type UserService interface {
	// Create creates new user and saves into store
	Create(userDTO *entity.UserDTO) (*entity.UserDTO, *errors.AppError)
	// FindByID returns user with matching id
	FindByID(userID string) (*entity.UserDTO, *errors.AppError)
	// FindByEmail returns user with matching email
	FindByEmail(email string) (*entity.UserDTO, *errors.AppError)
	// Update insert updated user values into store
	Update(userDTO *entity.UserDTO) (*entity.UserDTO, *errors.AppError)
	// Confirm verify user's account
	Confirm(userID string) *errors.AppError
	// Delete removes user entity from store
	Delete(id string) *errors.AppError
}

type userService struct {
	repository storage.UserRepository
}

func NewUserService(r storage.UserRepository) UserService {
	return &userService{
		repository: r,
	}
}

func (s *userService) Create(userDTO *entity.UserDTO) (*entity.UserDTO, *errors.AppError) {
	// validates userDTO fields
	validationError := entity.ValidateUser(userDTO)
	if validationError != "" {
		return nil, errors.BadRequest(validationError)
	}

	// hash user's password
	hashedPassword, err := entity.HashPassword(userDTO.Password)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	userDTO.Password = string(hashedPassword)

	createdUser, err := s.repository.Create(entity.ToUser(userDTO))

	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToUserDTO(createdUser), nil
}

func (s *userService) FindByEmail(email string) (*entity.UserDTO, *errors.AppError) {

	found, err := s.repository.FindByEmail(email)

	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	if found == nil {
		return nil, errors.NotFound("User is not found")
	}

	return entity.ToUserDTO(found), nil
}

func (s *userService) FindByID(userID string) (*entity.UserDTO, *errors.AppError) {

	found, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}
	if found == nil {
		return nil, errors.NotFound("User is not found")
	}

	return entity.ToUserDTO(found), nil
}

func (s *userService) Update(userDTO *entity.UserDTO) (*entity.UserDTO, *errors.AppError) {

	updated, err := s.repository.Update(entity.ToUser(userDTO))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToUserDTO(updated), nil
}

func (s *userService) Confirm(userID string) *errors.AppError {

	err := s.repository.Confirm(userID, true)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (s *userService) Delete(id string) *errors.AppError {

	err := s.repository.Delete(id)
	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}
