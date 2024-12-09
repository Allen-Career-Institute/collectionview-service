package service

import (
	"collectionview-service/internal/biz"
	"collectionview-service/internal/mongo"
	"collectionview-service/internal/utils"
	"context"
	"fmt"
	v1 "github.com/Allen-Career-Institute/common-protos/collection_view/v1"
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	pbrs "github.com/Allen-Career-Institute/common-protos/collection_view/v1/response"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

var ProviderSet = wire.NewSet(NewContentViewService)

type ContentViewService struct {
	v1.UnimplementedContentViewServiceServer
	mongoCollection mongo.MongoCollectionInterface
	log             *log.Helper
	bizHandler      *biz.CollectionBizHandler
}

func NewContentViewService(bizHandler *biz.CollectionBizHandler, mongoCollection mongo.MongoCollectionInterface) *ContentViewService {
	return &ContentViewService{
		UnimplementedContentViewServiceServer: v1.UnimplementedContentViewServiceServer{},
		mongoCollection:                       mongoCollection,
		log:                                   log.NewHelper(log.DefaultLogger),
		bizHandler:                            bizHandler,
	}
}

func (s *ContentViewService) GetCollectionView(ctx context.Context, req *pbrq.CollectionViewRequest) (*pbrs.CollectionViewResponse, error) {
	filter := bson.D{{"collection_id", req.CollectionId}}
	cursor, err := s.mongoCollection.List(ctx, filter, utils.Databasename, utils.LibCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to list from mongo collection: %w", err)
	}
	defer cursor.Close(ctx)

	var rawResults []bson.M
	if err = cursor.All(ctx, &rawResults); err != nil {
		return nil, fmt.Errorf("failed to decode cursor results: %w", err)
	}

	var results []*pbrs.CollectionView

	for _, rawResult := range rawResults {
		var Collection pbrs.CollectionView
		bsonBytes, err := bson.Marshal(rawResult)
		if err != nil {
			s.log.WithContext(ctx).Errorf("Failed to marshal BSON:", err)
			continue
		}
		err = bson.Unmarshal(bsonBytes, &Collection)
		if err != nil {
			s.log.WithContext(ctx).Errorf("Failed to unmarshal BSON:", err)
			continue
		}
		results = append(results, &Collection)
	}

	response := &pbrs.CollectionViewResponse{
		Collections: results,
	}

	return response, nil
}

func (s *ContentViewService) CreateCollection(ctx context.Context, req *pbrq.CreateCollectionViewRequest) (*pbrs.CreateCollectionViewResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, errors.New(http.StatusBadRequest, utils.InvalidRequestReason, utils.InvalidRequestMessage).WithMetadata(utils.GetErrorMetaData(err))
	}

	collectionID := utils.GenerateID(utils.LibraryPrefix)
	err := s.bizHandler.CreateCollection(ctx, req, collectionID)

	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to insert document into MongoDB: %v", err)
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return &pbrs.CreateCollectionViewResponse{
		CollectionId: collectionID,
		Message:      http.StatusText(201),
	}, nil
}

func (s *ContentViewService) UpdateCollection(ctx context.Context, req *pbrq.UpdateCollectionViewRequest) (*pbrs.UpdateCollectionViewResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, errors.New(http.StatusBadRequest, utils.InvalidRequestReason, utils.InvalidRequestMessage).WithMetadata(utils.GetErrorMetaData(err))
	}
	updateDBErr := s.bizHandler.UpdateCollection(ctx, req)

	if updateDBErr != nil {
		s.log.WithContext(ctx).Errorf("Failed to update MongoDB document: %v", updateDBErr)
		return nil, updateDBErr
	}
	response := &pbrs.UpdateCollectionViewResponse{
		CollectionId: req.CollectionId,
		Message:      http.StatusText(200),
	}
	return response, nil
}
