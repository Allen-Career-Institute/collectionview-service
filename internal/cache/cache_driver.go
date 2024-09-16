package cache

import (
	"collectionview-service/internal/conf"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type RedisStore struct {
	config   *conf.Redis
	DbClient redis.UniversalClient
	ctx      *context.Context
}

var (
	rs   *RedisStore
	once sync.Once
)

func NewRedisStore(data *conf.Redis) *RedisStore {
	redisConfig := data
	fmt.Println("connecting to redis")
	once.Do(func() {
		viper.SetConfigFile(redisConfig.CredFileLocation)
		viper.SetConfigType("json")
		err := viper.ReadInConfig()
		if err != nil {
			log.Errorf("Redis: Error reading secret file: %s", err)
		}
		viper.WatchConfig()

		addrs := []string{redisConfig.Addr}
		username := viper.GetString("username")
		password := viper.GetString("password")
		poolSize := redisConfig.PoolSize
		if poolSize == 0 {
			poolSize = 1
		}

		opts := &redis.ClusterOptions{
			Addrs:        addrs,
			Password:     password,
			PoolSize:     int(poolSize),
			ReadTimeout:  time.Duration(redisConfig.ReadTimeOutInMs) * time.Millisecond,
			DialTimeout:  time.Duration(redisConfig.DialTimeOutInMs) * time.Millisecond,
			WriteTimeout: time.Duration(redisConfig.WriteTimeOutInMs) * time.Millisecond,
		}

		if username != "" {
			opts.Username = username
		}

		if redisConfig.Tls {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		rdb := redis.NewClusterClient(opts)
		var ctx = context.Background()
		_, err = rdb.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("Error pinging Redis server: %v\n", err)
		} else {
			fmt.Println("successfully pinged redis")
		}
		rs = &RedisStore{
			config:   redisConfig,
			DbClient: rdb,
			ctx:      &ctx,
		}
	})
	fmt.Println("Connected to Redis")
	return rs
}
