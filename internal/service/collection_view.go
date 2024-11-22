package service

import (
	"collectionview-service/internal/cache"
	"collectionview-service/internal/mongo"
	"collectionview-service/internal/utils"
	"context"
	"fmt"
	v1 "github.com/Allen-Career-Institute/common-protos/collection_view/v1"
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	pbrs "github.com/Allen-Career-Institute/common-protos/collection_view/v1/response"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
)

var ProviderSet = wire.NewSet(NewContentViewService)

type ContentViewService struct {
	v1.UnimplementedContentViewServiceServer
	mongoCollection mongo.MongoCollectionInterface
	cacheCollection cache.CacheRepository
}

func NewContentViewService(mongoCollection mongo.MongoCollectionInterface, cacheCollection cache.CacheRepository) *ContentViewService {
	return &ContentViewService{
		UnimplementedContentViewServiceServer: v1.UnimplementedContentViewServiceServer{},
		mongoCollection:                       mongoCollection,
		cacheCollection:                       cacheCollection,
	}
}

// CreateCacheKey generates a cache key based on the fields of HomePageRequest
func CreateCacheKey(req *pbrq.CollectionViewRequest) string {
	// Use the relevant fields to generate the cache key
	return fmt.Sprintf("view:%s",
		req.CollectionId,
	)
}

// ListNcertBooks implements the ListNcertBooks method.
func (s *ContentViewService) GetCollectionView(ctx context.Context, req *pbrq.CollectionViewRequest) (*pbrs.CollectionViewResponse, error) {
	filter := bson.D{{"_id", req.CollectionId}}
	fmt.Println(filter, "filter")
	//cacheKey := CreateCacheKey(req)
	//Try to get the cached response
	//cachedResponse, err := s.cacheCollection.Get(ctx, cacheKey)
	//if err == nil {
	//	var response pbrs.CollectionViewResponse
	//	err = json.Unmarshal([]byte(cachedResponse), &response)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to unmarshal cached response: %v", err)
	//	}
	//}
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
			fmt.Println("Failed to marshal BSON:", err)
			continue
		}
		err = bson.Unmarshal(bsonBytes, &Collection)
		if err != nil {
			fmt.Println("Failed to unmarshal BSON:", err)
			continue
		}
		results = append(results, &Collection)
	}

	response := &pbrs.CollectionViewResponse{
		Collections: results,
	}

	//responseJSON, err := json.Marshal(response)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to marshal response to JSON: %v", err)
	//}
	////
	////// Define TTL for cache
	//ttl := utils.TTL

	//if err = s.cacheCollection.Set(ctx, cacheKey, string(responseJSON), ttl); err != nil {
	//	return nil, fmt.Errorf("failed to set response in cache: %w", err)
	//}

	//}
	return response, nil
}
