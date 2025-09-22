package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test koneksi Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️  Peringatan: Gagal menghubungkan ke Redis: %v", err)
		log.Println("⚠️  Aplikasi akan berjalan tanpa Redis, beberapa fitur mungkin terbatas")
		return nil
	}

	log.Printf("✅ Berhasil terhubung ke Redis: %s", pong)
	return rdb
}

// PingRedis melakukan ping test ke Redis untuk health check
func PingRedis(rdb *redis.Client) error {
	if rdb == nil {
		return fmt.Errorf("koneksi Redis tidak tersedia")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}

// GetRedisInfo mendapatkan informasi Redis untuk monitoring
func GetRedisInfo(rdb *redis.Client) (map[string]string, error) {
	if rdb == nil {
		return nil, fmt.Errorf("koneksi Redis tidak tersedia")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	info, err := rdb.Info(ctx, "server", "memory", "stats").Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	result["info"] = info

	// Dapatkan versi Redis
	if serverInfo, err := rdb.Info(ctx, "server").Result(); err == nil {
		result["server"] = serverInfo
	}

	return result, nil
}
