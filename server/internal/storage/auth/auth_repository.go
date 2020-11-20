package auth

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nozgurozturk/marvin/server/entity"
	"time"
)

type Repository struct {
	Client *redis.Client
}

// Creates new redis repository for auth
func NewRepository(client *redis.Client) *Repository {
	return &Repository{Client: client}
}

// Creates access token, uuid key value pair into redis db
func (r *Repository) CreateAuth(token *entity.Token) error {

	expires := time.Unix(token.Expires, 0)
	now := time.Now().UTC()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := r.Client.Set(ctx, token.Uuid, token.UserID, expires.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

// Gets access token from redis db
func (r *Repository) FindAuth(uuid string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := r.Client.Get(ctx, uuid).Result()
	if err != nil {
		return "", err
	}

	return id, nil
}

// Deletes access token from redis db
func (r *Repository) DeleteAuth(uuid string) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := r.Client.Del(ctx, uuid).Result()
	if err != nil {
		return err
	}

	return nil
}
