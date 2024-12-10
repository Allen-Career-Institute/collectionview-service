package data

import (
	"collectionview-service/internal/mongo"
	"collectionview-service/internal/utils"
	"context"
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type CollectionRepository struct {
	mongoCollection mongo.MongoCollectionInterface
	log             *log.Helper
}

type LearningMaterialRepo interface {
	CreateCollection(ctx context.Context, req *pbrq.CreateCollectionViewRequest, collectionID string) error
	UpdateCollection(ctx context.Context, req *pbrq.UpdateCollectionViewRequest) error
}

func NewCollectionRepository(mongoCollection mongo.MongoCollectionInterface, logger log.Logger) *CollectionRepository {
	return &CollectionRepository{
		mongoCollection: mongoCollection,
		log:             log.NewHelper(logger),
	}
}

func (repo *CollectionRepository) CreateCollection(ctx context.Context, innerpage *CollectionViewEntity) error {
	_, err := repo.mongoCollection.InsertDocument(ctx, utils.Databasename, utils.LibCollection, innerpage)
	return err
}

func (repo *CollectionRepository) UpdateCollection(ctx context.Context, req *pbrq.UpdateCollectionViewRequest) error {
	filter := bson.D{
		{utils.ID, req.CollectionId},
	}
	cursor, err := repo.mongoCollection.List(ctx, filter, utils.Databasename, utils.LibCollection)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("Failed to query MongoDB: %v", err)
		return err
	}
	defer cursor.Close(ctx)

	var rawResult bson.M
	if cursor.Next(ctx) {
		if err = cursor.Decode(&rawResult); err != nil {
			repo.log.WithContext(ctx).Errorf("Failed to decode cursor result: %v", err)
			return err
		}
	} else {
		repo.log.WithContext(ctx).Errorf("No document found matching the filter: %v", filter)
		return err
	}

	var cont CollectionViewEntity
	bsonBytes, err := bson.Marshal(rawResult)
	if err != nil {
		repo.log.WithContext(ctx).Errorf("Failed to marshal BSON result: %v", err)
		return err
	}
	if err = bson.Unmarshal(bsonBytes, &cont); err != nil {
		repo.log.WithContext(ctx).Errorf("Failed to unmarshal BSON: %v", err)
		return err
	}

	Mould(&cont, *req)
	updateDBErr := repo.mongoCollection.UpdateOne(ctx, utils.Databasename, utils.LibCollection, filter, cont)
	return updateDBErr
}
