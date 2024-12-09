package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionViewEntity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CollectionId string             `bson:"contentId,omitempty"`
	ViewDepth    int32              `bson:"class,omitempty"`
	ViewType     string             `bson:"createdAt,omitempty"`
	UpdatedAt    int64              `bson:"updatedAt,omitempty"`
	CreatedAt    int64              `bson:"updatedAt,omitempty"`
}
