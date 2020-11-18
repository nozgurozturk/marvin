package storage

import "github.com/nozgurozturk/marvin/notifier/entity"

type RepoRepository interface {
	Create(user *entity.Repo) (*entity.Repo, error)
	FindByUrlAndUserId(url string, userId string) (*entity.Repo, error)
	FindById(repoID string) (*entity.Repo, error)
	FindAll(userId string) ([]*entity.Repo, error)
	Delete(repoId string) error
	DeleteMany(userId string) error
}

type SubscriberRepository interface {
	Create(user *entity.Subscriber) (*entity.Subscriber, error)
	FindByEmailAndRepoId(email string, repoId string) (*entity.Subscriber, error)
	FindById(subID string) (*entity.Subscriber, error)
	Update(subscriber *entity.Subscriber) (*entity.Subscriber, error)
	Confirm(subId string, confirmed bool) error
	GetAll()([]*entity.Subscriber, error)
	FindAll(repoId string) ([]*entity.Subscriber, error)
	Delete(subId string) error
	DeleteMany(repoId string) error
}

