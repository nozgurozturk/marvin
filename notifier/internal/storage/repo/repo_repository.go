package repo

import (
	"context"
	"github.com/nozgurozturk/marvin/notifier/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("repositories")
	return &Repository{
		Collection: collection,
	}
}

func (r *Repository) Create(repo *entity.Repo) (*entity.Repo, error) {
	repo.CreatedAt = time.Now()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := r.Collection.InsertOne(ctx, &repo)
	if err != nil {
		return nil, err
	}
	repo.ID = result.InsertedID.(primitive.ObjectID)
	return repo, nil
}


func (r *Repository) FindById(repoID string) (*entity.Repo, error) {
	repo := new(entity.Repo)
	id, err := primitive.ObjectIDFromHex(repoID)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&repo)

	if err != nil {
		return nil, err
	}
	return repo, nil
}


func (r *Repository) FindByUrlAndUserId(url string, userId string) (*entity.Repo, error) {
	repo := new(entity.Repo)
	user, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = r.Collection.FindOne(ctx, bson.M{"path": url, "userId": user}).Decode(&repo)

	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) FindAll(userId string) ([]*entity.Repo, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	findFilter := bson.D{{"userId", bson.D{{"$in", bson.A{id}}}}}

	var repos []*entity.Repo

	findAllCursor, err := r.Collection.Find(ctx, findFilter)
	if findAllCursor != nil {
		if err = findAllCursor.All(ctx, &repos); err != nil {
			return nil, err
		}
	}

	return repos, nil
}


func (r *Repository) Delete(repoId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := primitive.ObjectIDFromHex(repoId)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteMany(userId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteMany(ctx, bson.D{{"userId", id}})
	if err != nil {
		return err
	}

	return nil
}