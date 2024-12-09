package mongo

import (
	"context"
	//document "github.com/Allen-Career-Institute/question-bank/internal/mongo/collection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionInterface interface {
	List(ctx context.Context, filter interface{}, dbName, collectionName string, opts ...*options.FindOptions) (*mongo.Cursor, error)
	InsertDocument(ctx context.Context, dbName, collectionName string, document interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, dbName, collectionName string, filter, updateData interface{}) error
}
