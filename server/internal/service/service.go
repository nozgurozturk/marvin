package service

import "github.com/nozgurozturk/marvin/server/internal/storage"

type Service interface {
	Auth() AuthService
	User() UserService
	Repo() RepoService
	Subscriber() SubscriberService
}

type service struct {
	auth       AuthService
	user       UserService
	repo       RepoService
	subscriber SubscriberService
}

func New(s storage.Store) *service {
	return &service{
		auth:       NewAuthService(s.Auths()),
		user:       NewUserService(s.Users()),
		repo:       NewRepoService(s.Repos()),
		subscriber: NewSubscriberService(s.Subscribers()),
	}
}

func (s *service) Auth() AuthService {
	return s.auth
}

func (s *service) User() UserService {
	return s.user
}

func (s *service) Repo() RepoService {
	return s.repo
}

func (s *service) Subscriber() SubscriberService {
	return s.subscriber
}
