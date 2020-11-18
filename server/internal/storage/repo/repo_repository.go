package repo

import (
	"context"
	"github.com/nozgurozturk/marvin/server/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

// Creates new mongo repository for git repositories
func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("repositories")
	return &Repository{
		Collection: collection,
	}
}

// Creates new git repository
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

// Finds git repository by id
func (r *Repository) FindByID(repoID string) (*entity.Repo, error) {

	repo := new(entity.Repo)

	id, err := primitive.ObjectIDFromHex(repoID)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&repo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

// Finds git repository by url and user id
func (r *Repository) FindByUrlAndUserID(url string, userID string) (*entity.Repo, error) {

	repo := new(entity.Repo)

	user, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = r.Collection.FindOne(ctx, bson.M{"path": url, "userID": user}).Decode(&repo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return repo, nil
}

// Finds git all repositories belongs to user
func (r *Repository) FindAll(userID string) ([]*entity.Repo, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	findFilter := bson.D{{"userID", bson.D{{"$in", bson.A{id}}}}}

	var repos []*entity.Repo

	findAllCursor, err := r.Collection.Find(ctx, findFilter)
	if findAllCursor != nil {
		if err = findAllCursor.All(ctx, &repos); err != nil {
			return nil, err
		}
	}

	return repos, nil
}

// Updates git repository's packages
func (r *Repository) UpdatePackages(repo *entity.Repo) (*entity.Repo, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := r.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", repo.ID}},
		bson.D{{"$set",
			bson.D{{"packageList", repo.PackageList}},
		}}).Decode(&repo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return repo, nil
}

// Deletes git repository
func (r *Repository) Delete(repoID string) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := primitive.ObjectIDFromHex(repoID)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}

// Deletes all git repositories belongs to user
// Run after deleting user
func (r *Repository) DeleteMany(userID string) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteMany(ctx, bson.D{{"userID", id}})
	if err != nil {
		return err
	}

	return nil
}
