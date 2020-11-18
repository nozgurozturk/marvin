package user

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

// Creates new mongo repository for users
func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("users")
	return &Repository{
		Collection: collection,
	}
}

// Creates new user
func (r *Repository) Create(user *entity.User) (*entity.User, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := r.Collection.InsertOne(ctx, &user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

// Finds user by id
func (r *Repository) FindByID(userID string) (*entity.User, error) {

	user := new(entity.User)

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Finds user by email
func (r *Repository) FindByEmail(email string) (*entity.User, error) {

	user := new(entity.User)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// Creates bson.D object with given values
// Using for user partial update
func createPartialUpdateOption(user *entity.User) (bson.D, error) {

	var updateOption bson.D

	if user.Password != "" {
		hashed, err := entity.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		updateOption = append(updateOption, bson.E{Key: "password", Value: hashed})
	}

	if user.Name != "" {
		updateOption = append(updateOption, bson.E{Key: "name", Value: user.Name})
	}

	if user.Email != "" {
		isValid := entity.ValidateEmail(user.Email)
		if isValid != "" {
			return nil, errors.New(isValid)
		}
		updateOption = append(updateOption, bson.E{Key: "email", Value: user.Email})
	}

	return updateOption, nil
}

// Updates user's fields
func (r *Repository) Update(user *entity.User) (*entity.User, error) {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	update, err := createPartialUpdateOption(user)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", user.ID}}, bson.D{{"$set", update}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

// Confirms user's email
// Sets true isConfirmed field
func (r *Repository) Confirm(userID string, confirmed bool) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := r.Collection.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"isConfirmed", confirmed}}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("User is not found")
	}

	return nil
}

// Deletes user from collection
func (r *Repository) Delete(userID string) error {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	return nil
}
