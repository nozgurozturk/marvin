package storage

type Store interface {
	Repos() RepoRepository
	Subscribers() SubscriberRepository
	Users() UserRepository
	Auths() AuthRepository
}
