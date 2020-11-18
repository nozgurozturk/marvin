package storage

type Store interface {
	Repos() RepoRepository
	Subscribers() SubscriberRepository
}
