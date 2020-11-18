package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Frequency string

const (
	Day  Frequency = "day"
	Week Frequency = "week"
	Hour Frequency = "hour"
)

type Notify struct {
	Hour      int          `json:"hour" bson:"hour"`
	Minute    int          `json:"minute" bson:"minute"`
	Weekday   time.Weekday `json:"weekday" bson:"weekday"`
	Frequency Frequency    `json:"frequency" bson:"frequency"`
}

type Subscriber struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RepoID      primitive.ObjectID `json:"repoID" bson:"repoID"`
	Email       string             `json:"email" bson:"email"`
	IsConfirmed bool               `json:"isConfirmed" bson:"isConfirmed"`
	Notify      *Notify            `json:"notify,omitempty" bson:"notify,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type SubscriberDTO struct {
	ID          *string `json:"id,omitempty"`
	RepoID      string  `json:"repoID"`
	Email       string  `json:"email"`
	IsConfirmed bool    `json:"isConfirmed"`
	Notify      *Notify `json:"notify,omitempty"`
}

type SubscriberIDRequest struct {
	ID string `json:"id"`
}

type SubscriberRequest struct {
	RepoID *string `json:"repoID,omitempty"`
	Email  *string `json:"email,omitempty"`
	Notify *Notify `json:"notify,omitempty"`
}

func ToSubscriberDTO(subscriber *Subscriber) *SubscriberDTO {

	id := subscriber.ID.Hex()

	return &SubscriberDTO{
		ID:          &id,
		RepoID:      subscriber.RepoID.Hex(),
		Email:       subscriber.Email,
		IsConfirmed: subscriber.IsConfirmed,
		Notify:      subscriber.Notify,
	}
}

func ToSubscriberDTOs(subscribers []*Subscriber) []*SubscriberDTO {

	subDTOs := make([]*SubscriberDTO, len(subscribers))

	for i, item := range subscribers {
		subDTOs[i] = ToSubscriberDTO(item)
	}

	return subDTOs
}

func ToSubscriber(subscribeDTO *SubscriberDTO) *Subscriber {

	repoID, _ := primitive.ObjectIDFromHex(subscribeDTO.RepoID)

	subscriber := &Subscriber{
		RepoID:      repoID,
		Email:       subscribeDTO.Email,
		IsConfirmed: subscribeDTO.IsConfirmed,
		Notify:      subscribeDTO.Notify,
	}

	if subscribeDTO.ID != nil {
		id, _ := primitive.ObjectIDFromHex(*subscribeDTO.ID)
		subscriber.ID = id
	} else {
		subscriber.ID = primitive.NilObjectID
	}

	return subscriber
}
