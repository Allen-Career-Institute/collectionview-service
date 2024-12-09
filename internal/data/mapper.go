package data

import (
	pbrq "github.com/Allen-Career-Institute/common-protos/collection_view/v1/request"
	"time"
)

func MouldReq(req *pbrq.CreateCollectionViewRequest, prefix string) *CollectionViewEntity {
	var collectionview CollectionViewEntity
	collectionview.CollectionId = prefix
	collectionview.CreatedAt = req.CreatedAt
	if req.CreatedAt == 0 {
		collectionview.CreatedAt = time.Now().Unix()
	}
	if req.UpdatedAt == 0 {
		collectionview.UpdatedAt = time.Now().Unix()
	}
	collectionview.UpdatedAt = req.UpdatedAt
	collectionview.ViewType = req.ViewType
	collectionview.ViewDepth = req.ViewDepth

	return &collectionview
}

func Mould(collection *CollectionViewEntity, req pbrq.UpdateCollectionViewRequest) {
	if req.ViewType != "" {
		collection.ViewType = req.ViewType
	}
	if req.ViewDepth != 0 {
		collection.ViewDepth = req.ViewDepth
	}

	collection.UpdatedAt = time.Now().Unix() // Always update `updated_at` to current time
}
