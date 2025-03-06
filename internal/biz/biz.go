package biz

import (
	"collectionview-service/internal/data"
	"collectionview-service/internal/mongo"
	"context"
	"fmt"
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewCollectionBizHandler)

type CollectionBizHandler struct {
	mongoCollection mongo.MongoCollectionInterface
	log             *log.Helper
	repo            *data.CollectionRepository
}

func NewCollectionBizHandler(mongoCollection mongo.MongoCollectionInterface, logger log.Logger) *CollectionBizHandler {
	return &CollectionBizHandler{
		mongoCollection: mongoCollection,
		log:             log.NewHelper(logger),
		repo:            data.NewCollectionRepository(mongoCollection, logger),
	}
}

func (s *CollectionBizHandler) CreateCollection(ctx context.Context, req *pbrq.CreateCollectionViewRequest, collectionID string) error {
	innerpage := data.MouldReq(req, collectionID)
	fmt.Println(innerpage, "shivansh")
	err := s.repo.CreateCollection(ctx, innerpage)
	return err
}

func (s *CollectionBizHandler) UpdateCollection(ctx context.Context, req *pbrq.UpdateCollectionViewRequest) error {
	updateDBErr := s.repo.UpdateCollection(ctx, req)
	return updateDBErr
}
