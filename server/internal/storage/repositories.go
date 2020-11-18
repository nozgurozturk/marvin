package storage

import "github.com/nozgurozturk/marvin/server/entity"

// UserRepository interface
type UserRepository interface {
	// Create insert entity to collection
	Create(user *entity.User) (*entity.User, error)
	// FindByID returns entity with matching id
	FindByID(userID string) (*entity.User, error)
	// FindByEmail returns entity with matching email
	FindByEmail(email string) (*entity.User, error)
	// Update insert updated entity values in collection
	Update(user *entity.User) (*entity.User, error)
	// Confirm updates entity when user verify email
	Confirm(userID string, confirmed bool) error
	// Delete removes entity from collection
	Delete(userID string) error
}

// RepoRepository interface
type RepoRepository interface {
	// Create insert entity to collection
	Create(user *entity.Repo) (*entity.Repo, error)
	// FindByID returns entity with matching id
	FindByID(repoID string) (*entity.Repo, error)
	// FindByEmail returns entity with matching url and user id
	FindByUrlAndUserID(url string, userID string) (*entity.Repo, error)
	// FindAll returns entities belongs to user
	FindAll(userID string) ([]*entity.Repo, error)
	// UpdatePackages insert updated packages into entity
	UpdatePackages(repo *entity.Repo) (*entity.Repo, error)
	// Delete removes entity from collection
	Delete(repoID string) error
	// Delete removes all entities belongs to user
	DeleteMany(userID string) error
}

// SubscriberRepository interface
type SubscriberRepository interface {
	// Create insert entity to collection
	Create(user *entity.Subscriber) (*entity.Subscriber, error)
	// FindByEmail returns entity with matching email and repository id
	FindByEmailAndRepoID(email string, repoID string) (*entity.Subscriber, error)
	// FindByID returns entity with matching id
	FindByID(subID string) (*entity.Subscriber, error)
	// FindAll returns entities belongs to git repository
	FindAll(repoID string) ([]*entity.Subscriber, error)
	// Update insert updated entity values in collection
	Update(subscriber *entity.Subscriber) (*entity.Subscriber, error)
	// Confirm updates entity when subscriber verify email
	Confirm(subID string, confirmed bool) error
	// GetAll returns all entities in collection
	GetAll()([]*entity.Subscriber, error)
	// Delete removes entity from collection
	Delete(subID string) error
	// Delete removes all entities belongs to git repository
	DeleteMany(repoID string) error
}

// AuthRepository interface
type AuthRepository interface {
	// CreateAuth insert entity to store
	CreateAuth(token *entity.Token) error
	// FindAuth returns entity with matching uuid
	FindAuth(uuid string) (string, error)
	// Delete removes entity from store
	DeleteAuth(uuid string) error
}
