package mongo

import (
	"context"
	//document "github.com/Allen-Career-Institute/question-bank/internal/mongo/collection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionInterface interface {
	//InsertDocument(context.Context, interface{}, string, string) (interface{}, error)
	//InsertDocuments(context.Context, []interface{}, string, string) ([]interface{}, error)
	Get(ctx context.Context, filter interface{}, dbName, collectionName string) (*mongo.SingleResult, error)
	List(ctx context.Context, filter interface{}, dbName, collectionName string, opts ...*options.FindOptions) (*mongo.Cursor, error)
	//Aggregate(ctx context.Context, pipeline mongo.Pipeline, dbName, collectionName string) (*mongo.Cursor, error)
	//Update(ctx context.Context, filter, updateData interface{}, dbName, collectionName string) error
	//BatchGet(ctx context.Context, filter interface{}, dbName, collectionName string) (*mongo.Cursor, error)
	//UpsertDynamicTags(ctx context.Context, tags *document.QuestionDynamicTags) error
	//GetDynamicTags(ctx context.Context, questionIDs []string) ([]*document.QuestionDynamicTags, error)
	//GetQuestionsByOldQuestionID(ctx context.Context, questionIds []int64) ([]*document.QuestionDocument, error)
}
