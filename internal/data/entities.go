package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollectionViewEntity struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CollectionId string             `bson:"collectionId,omitempty"`
	ViewDepth    int32              `bson:"ViewDepth,omitempty"`
	ViewType     string             `bson:"ViewType,omitempty"`
	UpdatedAt    int64              `bson:"updatedAt,omitempty"`
	CreatedAt    int64              `bson:"createdAt,omitempty"`
}
