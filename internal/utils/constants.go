package utils

import "time"

const (
	MongoErrorTag         = "MONGO_ERROR"
	Databasename          = "collection-view"
	LibCollection         = "collection-repo"
	TTL                   = 36 * time.Hour
	defaultSize           = 21
	defaultAlphabet       = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	InvalidRequestReason  = "INVALID_REQUEST"
	InvalidRequestMessage = "Invalid Request"
	Collection            = "collection"
	LibraryPrefix         = "lib"
	ID                    = "collectionId"
	IsActive              = "isActive"
	CreatedAt             = "createdAt"
	UpdatedAt             = "updatedAt"
)
