package mongo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	START_TIME_LOG = "Start Time to Fetch Doc in mongo = %v"
)

type mongoCollectionImpl struct {
	client *mongo.Client
	log    *log.Helper
}

func (q *mongoCollectionImpl) Get(ctx context.Context, filter interface{}, dbName, collectionName string) (*mongo.SingleResult, error) {
	// Start a new session for the transaction
	q.log.WithContext(ctx).Debugf(START_TIME_LOG, time.Now())
	session, err := q.client.StartSession()
	defer session.EndSession(ctx)
	// Start the transaction
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	// Get a handle to the database and collection
	db := q.client.Database(dbName)
	collection := db.Collection(collectionName)

	// Query the collection
	result := collection.FindOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("No document found")
			return nil, nil
		}
		// Abort the transaction in case of an error
		err = session.AbortTransaction(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (q *mongoCollectionImpl) List(ctx context.Context, filter interface{}, dbName, collectionName string, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Log the start time of the query
	q.log.WithContext(ctx).Debugf(START_TIME_LOG, time.Now())

	// Get a handle to the database and collection
	db := q.client.Database(dbName)
	collection := db.Collection(collectionName)

	// Query the collection
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		q.log.WithContext(ctx).Errorf("Error during query execution: %v", err)
		return nil, err
	}

	// Check for any cursor errors
	if err = cursor.Err(); err != nil {
		q.log.WithContext(ctx).Errorf("Error during cursor iteration: %v", err)
	}

	return cursor, err
}

func (lr *mongoCollectionImpl) InsertDocument(ctx context.Context, dbName, collectionName string, documents interface{}) (interface{}, error) {
	db := lr.client.Database(dbName)
	collection := db.Collection(collectionName)

	data, err := collection.InsertOne(ctx, documents)
	if err != nil {
		return nil, err
	}

	return data.InsertedID, nil
}

func (lr *mongoCollectionImpl) UpdateOne(ctx context.Context, dbName, collectionName string, filter, updateData interface{}) error {
	db := lr.client.Database(dbName)
	collection := db.Collection(collectionName)

	update := bson.M{"$set": updateData}
	_, err := collection.UpdateOne(ctx, filter, update)

	return err
}

func NewMongoCollectionImpl(data *Data, logger *log.Helper) MongoCollectionInterface {
	return &mongoCollectionImpl{data.client, logger}
}
