package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"github.com/nozgurozturk/marvin/server/internal/storage/auth"
	"github.com/nozgurozturk/marvin/server/internal/storage/repo"
	"github.com/nozgurozturk/marvin/server/internal/storage/subscriber"
	"github.com/nozgurozturk/marvin/server/internal/storage/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DB struct {
	mongo       *mongo.Database
	redis       *redis.Client
	auths       AuthRepository
	repos       RepoRepository
	users       UserRepository
	subscribers SubscriberRepository
}
// Connects MongoDB and returns mongo.Database struct
func MongoConnect() (*mongo.Database, error) {
	cnf := config.Get().Mongo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If you use mongodb atlas for development
	// connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s", cnf.Username, cnf.Password, cnf.Host, cnf.DBName, cnf.Query)

	// If you use local development
	connectionString := fmt.Sprintf("mongodb://%s/%s", cnf.Host, cnf.DBName)

	// If you use local development with auth
	// connectionString := fmt.Sprintf("mongodb://%s:%s@%s/%s?%s", cnf.Username, cnf.Password, cnf.Host, cnf.DBName, cnf.Query)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Mongo is successfully connected to Application")

	db := client.Database(cnf.DBName)

	return db, err
}

// Connects Redis and returns redis.Client struct
func RedisConnect() (*redis.Client, error) {
	cnf := config.Get().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cnf.Address,
		Username: cnf.UserName,
		Password: cnf.Password,
		DB:       cnf.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Redis is successfully connected to Application")
	return client, nil
}

// Creates new database layer
func New(mongo *mongo.Database, redis *redis.Client) *DB {
	return &DB{
		mongo:       mongo,
		redis:       redis,
		repos:       repo.NewRepository(mongo),
		users:       user.NewRepository(mongo),
		subscribers: subscriber.NewRepository(mongo),
		auths:       auth.NewRepository(redis),
	}
}

// Returns git repository mongo repository
func (db *DB) Repos() RepoRepository {
	return db.repos
}

// Returns user mongo repository
func (db *DB) Users() UserRepository {
	return db.users
}

// Returns subscriber mongo repository
func (db *DB) Subscribers() SubscriberRepository {
	return db.subscribers
}

// Returns auth redis repository
func (db *DB) Auths() AuthRepository {
	return db.auths
}
