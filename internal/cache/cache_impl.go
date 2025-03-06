package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type cacheImpl struct {
	ctx        *context.Context
	cacheStore *RedisStore
}

func (c *cacheImpl) DeleteByPattern(ctx context.Context, pattern string) error {
	cursor := uint64(0)
	keysToDelete := make([]string, 0)

	for {
		keys, nextCursor, err := c.cacheStore.DbClient.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return err
		}
		keysToDelete = append(keysToDelete, keys...)
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	if len(keysToDelete) > 0 {
		if err := c.cacheStore.DbClient.Del(ctx, keysToDelete...).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (c *cacheImpl) Delete(ctx context.Context, key string) (int64, error) {
	sc := c.cacheStore.DbClient.Del(ctx, key)
	return sc.Result()
}

func (c *cacheImpl) SetMultiple(ctx context.Context, dataMap map[string]interface{}, ttl time.Duration) (int, error) {
	if len(dataMap) == 0 {
		return 0, nil
	}
	pipe := c.cacheStore.DbClient.TxPipeline()
	stringData := make([]interface{}, 0, len(dataMap)*2)
	for key, value := range dataMap {
		jsonString, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}
		stringData = append(stringData, key, string(jsonString))
	}

	if err := pipe.MSet(ctx, stringData...).Err(); err != nil {
		return 0, err
	}

	for key := range dataMap {
		pipe.Expire(ctx, key, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return len(dataMap), nil
}
func (c *cacheImpl) GetByPattern(ctx context.Context, pattern string, count int) ([]string, error) {
	var keys []string
	cursor := uint64(0)

	for {
		result, nextCursor, err := c.cacheStore.DbClient.Scan(ctx, cursor, pattern, int64(count)).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, result...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	values, err := c.cacheStore.DbClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var stringValues []string
	for _, val := range values {
		if val != nil {
			stringValues = append(stringValues, val.(string))
		} else {
			stringValues = append(stringValues, "")
		}
	}

	return stringValues, nil

}

func (c *cacheImpl) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := c.cacheStore.DbClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *cacheImpl) Get(ctx context.Context, key string) (string, error) {
	sc := c.cacheStore.DbClient.Get(ctx, key)
	return sc.Result()
}

func NewCacheImpl(cache *RedisStore) CacheRepository {
	return &cacheImpl{cacheStore: cache}
}
