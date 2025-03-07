package service

import (
	"collectionview-service/internal/biz"
	"collectionview-service/internal/cache"
	"collectionview-service/internal/mongo"
	"collectionview-service/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	v1 "github.com/Allen-Career-Institute/common-protos/collection_view/v1"
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	pbrs "github.com/Allen-Career-Institute/common-protos/collection_view/v1/response"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProviderSet = wire.NewSet(NewContentViewService)

type ContentViewService struct {
	v1.UnimplementedContentViewServiceServer
	mongoCollection mongo.MongoCollectionInterface
	log             *log.Helper
	bizHandler      *biz.CollectionBizHandler
	CacheStore      cache.CacheRepository
}

func NewContentViewService(bizHandler *biz.CollectionBizHandler, mongoCollection mongo.MongoCollectionInterface, cacheStore cache.CacheRepository) *ContentViewService {
	return &ContentViewService{
		UnimplementedContentViewServiceServer: v1.UnimplementedContentViewServiceServer{},
		mongoCollection:                       mongoCollection,
		log:                                   log.NewHelper(log.DefaultLogger),
		bizHandler:                            bizHandler,
		CacheStore:                            cacheStore,
	}
}

func (s *ContentViewService) GetCollectionView(ctx context.Context, req *pbrq.CollectionViewRequest) (*pbrs.CollectionViewResponse, error) {
	filter := bson.D{{utils.ID, req.CollectionId}}
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

func (s *ContentViewService) CreateCollectionView(ctx context.Context, req *pbrq.CreateCollectionViewRequest) (*pbrs.CreateCollectionViewResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, errors.New(http.StatusBadRequest, utils.InvalidRequestReason, utils.InvalidRequestMessage).WithMetadata(utils.GetErrorMetaData(err))
	}

	collectionID := req.CollectionId
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

func (s *ContentViewService) UpdateCollectionView(ctx context.Context, req *pbrq.UpdateCollectionViewRequest) (*pbrs.UpdateCollectionViewResponse, error) {
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

type ReelData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TopicId   string             `bson:"topicId"`
	SubjectId string             `bson:"subjectId"`
	VideoID   string             `bson:"videoID"`
	URL       string             `bson:"url"`
	Title     string             `bson:"title"`
	Subtitle  string             `bson:"subtitle"`
	Taxonomy  string             `bson:"taxonomy"`
}

func (s *ContentViewService) GetReelCollection(ctx context.Context, req *pbrq.GetReelCollectionRequest) (*pbrs.GetReelCollectionResponse, error) {
	const requiredCount = 5
	var allVideos []ReelData
	userID := req.UserId
	cacheKey := fmt.Sprintf("watched_reels:%s", userID)
	fmt.Println(cacheKey)
	// 1. Get already watched reels from cache
	watchedReelsJSON, err := s.CacheStore.Get(ctx, cacheKey)
	var watchedReels map[string]struct{}
	if err == nil && watchedReelsJSON != "" {
		_ = json.Unmarshal([]byte(watchedReelsJSON), &watchedReels)
	} else {
		watchedReels = make(map[string]struct{})
	}
	fmt.Println("watchedReels", watchedReels)
	// 2. Fetch unwatched videos by Topic
	topicFilter := bson.M{"topicId": req.TopicId, "videoID": bson.M{"$nin": getKeys(watchedReels)}}
	topicVideos, err := s.fetchVideos(ctx, topicFilter, 30)
	if err != nil {
		return nil, err
	}
	if len(topicVideos) >= requiredCount {
		return s.prepareResponse(ctx, topicVideos[:requiredCount], watchedReels, cacheKey)
	}
	allVideos = append(allVideos, topicVideos...)

	// 3. Fallback to Subject filter if not enough
	if len(allVideos) < requiredCount {
		subjectFilter := bson.M{"subjectId": req.SubjectId, "videoID": bson.M{"$nin": getKeys(watchedReels)}}
		subjectVideos, err := s.fetchVideos(ctx, subjectFilter, 30)
		if err != nil {
			return nil, err
		}
		allVideos = append(allVideos, subjectVideos...)
		if len(allVideos) >= requiredCount {
			return s.prepareResponse(ctx, allVideos[:requiredCount], watchedReels, cacheKey)
		}
	}

	// 4. Fallback to random videos if still not enough
	if len(allVideos) < requiredCount {
		randomFilter := bson.M{"videoID": bson.M{"$nin": getKeys(watchedReels)}}
		randomVideos, err := s.fetchVideos(ctx, randomFilter, 30)
		if err != nil {
			return nil, err
		}
		allVideos = append(allVideos, randomVideos...)
		if len(allVideos) >= requiredCount {
			return s.prepareResponse(ctx, allVideos[:requiredCount], watchedReels, cacheKey)
		}
	}

	// 5. If still not enough, fill with random videos
	if len(allVideos) < requiredCount {
		needed := requiredCount - len(allVideos)
		randomFillVideos, err := s.fetchVideos(ctx, bson.M{}, int64(needed))
		if err != nil {
			return nil, err
		}
		allVideos = append(allVideos, randomFillVideos...)
	}

	return s.prepareResponse(ctx, allVideos, watchedReels, cacheKey)
}

// fetchVideos fetches videos from MongoDB based on a filter and limit
func (s *ContentViewService) fetchVideos(ctx context.Context, filter bson.M, limit int64) ([]ReelData, error) {
	opts := options.Find().SetLimit(limit)
	cursor, err := s.mongoCollection.List(ctx, filter, utils.Databasename, utils.LibCollection, opts)
	if err != nil {
		return nil, err
	}
	var videos []ReelData
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(videos), func(i, j int) { videos[i], videos[j] = videos[j], videos[i] })
	return videos, nil
}

// prepareResponse prepares the final response and updates cache
func (s *ContentViewService) prepareResponse(ctx context.Context, videos []ReelData, watchedReels map[string]struct{}, cacheKey string) (*pbrs.GetReelCollectionResponse, error) {
	response := &pbrs.GetReelCollectionResponse{
		Reels: make([]*pbrs.ReelData, len(videos)),
	}
	for i, video := range videos {
		response.Reels[i] = &pbrs.ReelData{
			Id:        video.VideoID,
			Url:       video.URL,
			Title:     video.Title,
			Subtitle:  video.Subtitle,
			SubjectId: video.SubjectId,
			TopicId:   video.TopicId,
			Taxonomy:  video.Taxonomy,
		}
		watchedReels[video.VideoID] = struct{}{}
	}
	// Update cache with new watched reels
	updatedCache, _ := json.Marshal(watchedReels)
	_ = s.CacheStore.Set(ctx, cacheKey, string(updatedCache), 24*time.Hour)
	return response, nil
}

// getKeys returns a list of keys from a map
func getKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
