package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) RedisRepository {
	return &redisRepository{rdb: rdb}
}

func (r *redisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepository) Get(key string) (string, error) {
	if r.rdb == nil {
		return "", redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (r *redisRepository) GetJSON(key string, dest interface{}) error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

func (r *redisRepository) SetJSON(key string, value interface{}, ttl time.Duration) error {
	if r.rdb == nil {
		return redis.Nil
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Set(ctx, key, jsonData, ttl).Err()
}

func (r *redisRepository) Delete(key string) error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Del(ctx, key).Err()
}

func (r *redisRepository) Exists(key string) (bool, error) {
	if r.rdb == nil {
		return false, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.rdb.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *redisRepository) Increment(key string) (int64, error) {
	if r.rdb == nil {
		return 0, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Incr(ctx, key).Result()
}

func (r *redisRepository) Decrement(key string) (int64, error) {
	if r.rdb == nil {
		return 0, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Decr(ctx, key).Result()
}

func (r *redisRepository) SetExpire(key string, ttl time.Duration) error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.Expire(ctx, key, ttl).Err()
}

func (r *redisRepository) GetTTL(key string) (time.Duration, error) {
	if r.rdb == nil {
		return 0, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.TTL(ctx, key).Result()
}

func (r *redisRepository) GetKeys(pattern string) ([]string, error) {
	if r.rdb == nil {
		return nil, redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (r *redisRepository) FlushAll() error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.rdb.FlushAll(ctx).Err()
}

func (r *redisRepository) Ping() error {
	if r.rdb == nil {
		return redis.Nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.rdb.Ping(ctx).Result()
	return err
}
