package utils

import "time"

const (
	MongoErrorTag          = "MONGO_ERROR"
	Databasename           = "collection-view"
	LibCollection          = "collection-repo"
	TTL                    = 36 * time.Hour
	DefaultCacheEntryCount = 5
)
