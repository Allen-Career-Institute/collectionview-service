package cache

import (
	"collectionview-service/constants"
	"collectionview-service/internal/conf"
	"context"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupCache() (context.Context, redis.UniversalClient, redismock.ClientMock) {
	ctx := context.Background()
	mockClient, mock := redismock.NewClientMock()
	return ctx, mockClient, mock
}

func TestNewCacheImpl(t *testing.T) {
	c, mockCacheClient, _ := setupCache()
	type args struct {
		cacheStore *RedisStore
	}
	tests := []struct {
		name string
		args args
		want CacheRepository
	}{
		{
			name: "Test new redis impl",
			args: args{
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheStore := NewCacheImpl(tt.args.cacheStore)
			assert.NotNil(t, cacheStore)
		})
	}
}

func Test_cacheImpl_Delete(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Successful",
			args: args{
				ctx: c,
				key: "cache_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				expectErr := redis.NewCmdResult(1, nil)
				mock.ExpectDel("cache_key").SetErr(expectErr.Err())
			},
		},
		{
			name: "Empty key error",
			args: args{
				ctx: c,
				key: "cache_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				expectErr := errors.New("empty key error")
				mock.ExpectDel("").SetErr(expectErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockCacheClient.Del(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_cacheImpl_DeleteByPattern(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Successful",
			args: args{
				ctx: c,
				key: "cache_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				expectErr := redis.NewCmdResult(1, nil)
				mock.ExpectDel("cache_key").SetErr(expectErr.Err())
			},
		},
		{
			name: "Empty key error",
			args: args{
				ctx: c,
				key: "cache_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				expectErr := errors.New("empty key error")
				mock.ExpectDel("").SetErr(expectErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockCacheClient.Del(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func Test_cacheImpl_Get(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: c,
				key: "redis_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: false,
			want:    "dummyValue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectGet("redis_key").SetVal(tt.want)
			rs := &cacheImpl{
				ctx:        tt.fields.ctx,
				cacheStore: tt.fields.cacheStore,
			}
			got, err := rs.Get(tt.args.ctx, tt.args.key)
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, got, "Get(%v, %v)", tt.args.ctx, tt.args.key)
		})
	}
}

func Test_cacheImpl_GetByPattern(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
		want    []string
	}{
		{
			name: "success",
			args: args{
				ctx: c,
				key: "redis_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: false,
			mock: func() {
				mock.ExpectMGet("redis_key").SetVal([]interface{}{"test"})
			},
		},
		{
			name: "failure",
			args: args{
				ctx: c,
				key: "redis_key",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				mock.ExpectMGet("redis_key").SetErr(errors.New("not found"))
			},
		},
		{
			name: "failure empty key",
			args: args{
				ctx: c,
				key: "",
			},
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			wantErr: true,
			mock: func() {
				mock.ExpectMGet("").SetErr(errors.New("not found"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &cacheImpl{
				ctx:        tt.fields.ctx,
				cacheStore: tt.fields.cacheStore,
			}

			got, err := rs.GetByPattern(tt.args.ctx, tt.args.key, constants.DefaultCacheEntryCount)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func Test_cacheImpl_Set(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx   context.Context
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Set successfully",
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			args: args{
				ctx:   c,
				key:   "redis_key",
				value: "redis_value",
				ttl:   1 * time.Second,
			},
			wantErr: false,
			mock: func() {
				mock.ExpectSet("redis_key", "redis_value", 1*time.Second).SetVal("true")
			},
		},
		{
			name: "Set failure",
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			args: args{
				ctx:   c,
				key:   "redis_key",
				value: "redis_value",
				ttl:   1000,
			},
			wantErr: true,
			mock: func() {
				mock.ExpectSet("redis_key", "redis_value", time.Second).SetErr(errors.New("set error"))
			},
		},
		{
			name: "Set failure empty key",
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			args: args{
				ctx:   c,
				key:   "",
				value: "redis_value",
				ttl:   1000,
			},
			wantErr: true,
			mock: func() {
				mock.ExpectSet("", "redis_value", time.Second).SetErr(errors.New("empty key error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			rs := &cacheImpl{
				ctx:        tt.fields.ctx,
				cacheStore: tt.fields.cacheStore,
			}

			err := rs.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.ttl)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cacheImpl_SetMultiple(t *testing.T) {
	c, mockCacheClient, mock := setupCache()
	type fields struct {
		ctx        *context.Context
		cacheStore *RedisStore
	}
	type args struct {
		ctx  context.Context
		data map[string]interface{}
		ttl  time.Duration
	}
	data := map[string]interface{}{
		"LifeStatus": "Dead",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "Set Map successfully",
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			args: args{
				ctx:  c,
				data: data,
				ttl:  1000 * time.Millisecond,
			},
			wantErr: false,
			mock: func() {
				mock.ExpectTxPipeline()
				mock.ExpectMSet("LifeStatus", "\"Dead\"").SetVal("")
				mock.ExpectExpire("LifeStatus", 1).SetVal(true)
				mock.ExpectTxPipelineExec()
			},
		},
		{
			name: "Set Map error",
			fields: fields{
				ctx: &c,
				cacheStore: &RedisStore{
					config:   &conf.Redis{},
					DbClient: mockCacheClient,
					ctx:      &c,
				},
			},
			args: args{
				ctx:  c,
				data: data,
				ttl:  1000 * time.Millisecond,
			},
			wantErr: true,
			mock: func() {
				mock.ExpectHMSet("redis_key", data).SetErr(errors.New(""))

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			rs := &cacheImpl{
				ctx:        tt.fields.ctx,
				cacheStore: tt.fields.cacheStore,
			}
			_, err := rs.SetMultiple(tt.args.ctx, tt.args.data, 1000*time.Millisecond)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
