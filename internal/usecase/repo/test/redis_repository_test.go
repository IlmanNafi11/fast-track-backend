package repo_test

import (
	"fiber-boiler-plate/internal/usecase/repo"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_SetAndGet(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	err := redisRepo.Set("test-key", "test-value", time.Minute)
	assert.Equal(t, redis.Nil, err)

	value, err := redisRepo.Get("test-key")
	assert.Equal(t, "", value)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_SetJSON(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	testData := map[string]interface{}{
		"key":    "value",
		"number": 123,
	}

	err := redisRepo.SetJSON("test-json", testData, time.Hour)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_GetJSON(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	var result map[string]interface{}
	err := redisRepo.GetJSON("test-json", &result)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_Delete(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	err := redisRepo.Delete("test-key")
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_Exists(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	exists, err := redisRepo.Exists("test-key")
	assert.False(t, exists)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_Increment(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	result, err := redisRepo.Increment("counter")
	assert.Equal(t, int64(0), result)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_Decrement(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	result, err := redisRepo.Decrement("counter")
	assert.Equal(t, int64(0), result)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_SetExpire(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	err := redisRepo.SetExpire("test-key", time.Minute)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_GetTTL(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	ttl, err := redisRepo.GetTTL("test-key")
	assert.Equal(t, time.Duration(0), ttl)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_FlushAll(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	err := redisRepo.FlushAll()
	assert.Equal(t, redis.Nil, err)
}

func TestRedisRepository_GetKeys(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	keys, err := redisRepo.GetKeys("test-pattern*")
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, keys)
}

func TestRedisRepository_Ping(t *testing.T) {
	redisRepo := repo.NewRedisRepository(nil)

	err := redisRepo.Ping()
	assert.Equal(t, redis.Nil, err)
}
