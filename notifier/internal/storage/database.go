package storage

import (
	"context"
	"fmt"
	"github.com/nozgurozturk/marvin/notifier/internal/config"
	"github.com/nozgurozturk/marvin/notifier/internal/storage/repo"
	"github.com/nozgurozturk/marvin/notifier/internal/storage/subscriber"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DB struct {
	mongo       *mongo.Database
	repos       RepoRepository
	subscribers SubscriberRepository
}

func MongoConnect() (*mongo.Database, error) {
	cnf := config.Get().Mongo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("mongodb://%s/%s", cnf.Host, cnf.DBName)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Mongo is successfully connected to Application")

	db := client.Database(cnf.DBName)

	return db, err
}

func New(mongo *mongo.Database) *DB {
	return &DB{
		mongo:       mongo,
		repos:       repo.NewRepository(mongo),
		subscribers: subscriber.NewRepository(mongo),
	}
}

func (db *DB) Repos() RepoRepository {
	return db.repos
}

func (db *DB) Subscribers() SubscriberRepository {
	return db.subscribers
}

