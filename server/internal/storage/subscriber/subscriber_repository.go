package subscriber

import (
	"context"
	"errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

// Creates new mongo repository for subscribers
func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("subscribers")
	return &Repository{
		Collection: collection,
	}
}

// Creates subscriber
func (r *Repository) Create(sub *entity.Subscriber) (*entity.Subscriber, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := r.Collection.InsertOne(ctx, &sub)
	if err != nil {
		return nil, err
	}

	sub.ID = result.InsertedID.(primitive.ObjectID)

	return sub, nil
}

// Finds subscriber by id
func (r *Repository) FindByID(subID string) (*entity.Subscriber, error) {

	sub := new(entity.Subscriber)

	id, err := primitive.ObjectIDFromHex(subID)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return sub, nil
}

// Finds subscriber by email and git repository id
func (r *Repository) FindByEmailAndRepoID(email string, repoID string) (*entity.Subscriber, error) {

	sub := new(entity.Subscriber)

	repo, err := primitive.ObjectIDFromHex(repoID)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = r.Collection.FindOne(ctx, bson.M{"email": email, "repoID": repo}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return sub, nil
}

// Finds all subscribers belongs to git repository
func (r *Repository) FindAll(repoID string) ([]*entity.Subscriber, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := primitive.ObjectIDFromHex(repoID)
	if err != nil {
		return nil, err
	}

	findFilter := bson.D{{"repoID", bson.D{{"$in", bson.A{id}}}}}
	findAllCursor, err := r.Collection.Find(ctx, findFilter)
	if err != nil {
		return nil, err
	}

	var subscribers []*entity.Subscriber
	if findAllCursor != nil {
		if err = findAllCursor.All(ctx, &subscribers); err != nil {
			return nil, err
		}
	}

	return subscribers, nil
}

// Gets all subscribers in collection
func (r *Repository) GetAll() ([]*entity.Subscriber, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var subscribers []*entity.Subscriber

	findAllCursor, err := r.Collection.Find(ctx, bson.D{})
	if findAllCursor != nil {
		if err = findAllCursor.All(ctx, &subscribers); err != nil {
			return nil, err
		}
	}

	return subscribers, nil
}

// Creates bson.D object with given values
// Using for git repository partial update
func createPartialUpdateOption(subscriber *entity.Subscriber) (bson.D, error) {

	var updateOption bson.D

	if subscriber.IsConfirmed {
		updateOption = append(updateOption, bson.E{Key: "isConfirmed", Value: subscriber.IsConfirmed})
	}

	if subscriber.Email != "" {
		isValid := entity.ValidateEmail(subscriber.Email)
		if isValid != "" {
			return nil, errors.New(isValid)
		}
		updateOption = append(updateOption, bson.E{Key: "email", Value: subscriber.Email})
	}

	if subscriber.Notify != nil {
		updateOption = append(updateOption, bson.E{Key: "notify", Value: subscriber.Notify})
	}

	return updateOption, nil
}

// Updates subscriber's fields
func (r *Repository) Update(subscriber *entity.Subscriber) (*entity.Subscriber, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	update, err := createPartialUpdateOption(subscriber)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", subscriber.ID}}, bson.D{{"$set", update}}).Decode(&subscriber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return subscriber, nil
}

// Confirms subscribers's email
// Sets true isConfirmed field
func (r *Repository) Confirm(subID string, confirmed bool) error {

	id, err := primitive.ObjectIDFromHex(subID)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := r.Collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"isConfirmed", confirmed}}}})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("Subscriber is not found")
	}

	return nil
}

// Deletes subscriber from collection
func (r *Repository) Delete(subID string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := primitive.ObjectIDFromHex(subID)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}

// Deletes subscribers belongs to git repository
// Run after deleting git repository
func (r *Repository) DeleteMany(repoId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	id, err := primitive.ObjectIDFromHex(repoId)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteMany(ctx, bson.D{{"repoId", id}})
	if err != nil {
		return err
	}

	return nil
}
