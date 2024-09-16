package cache

import (
	"collectionview-service/internal/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRedisUniversalClient struct {
	mock.Mock
}
type MockConfDataRedis struct {
	mock.Mock
}

func (m *MockConfDataRedis) GetAddr() string {
	args := m.Called()
	return args.String(0)
}

func TestNewRedisStore(t *testing.T) {
	mockConfRedis := new(MockConfDataRedis)
	mockConfRedis.On("GetAddr").Return("localhost:6379")

	mockRedisClient := new(MockRedisUniversalClient)
	mockRedisClient.On("Close").Return(nil)

	redisConfig := &conf.Redis{
		Addr: mockConfRedis.GetAddr(),
	}
	rs := NewRedisStore(redisConfig)

	assert.NotNil(t, rs)
	assert.Equal(t, redisConfig, rs.config)
	assert.NotNil(t, rs.ctx)
	mockConfRedis.AssertExpectations(t)

}
