package config

import (
	"context"
	"fiber-boiler-plate/internal/helper"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		helper.Warn("Peringatan: Gagal menghubungkan ke Redis", logrus.Fields{
			"error": err.Error(),
		})
		helper.Warn("Aplikasi akan berjalan tanpa Redis, beberapa fitur mungkin terbatas")
		return nil
	}

	helper.Info("Berhasil terhubung ke Redis", logrus.Fields{
		"response": pong,
		"host":     cfg.Redis.Host,
		"port":     cfg.Redis.Port,
	})
	return rdb
}

func PingRedis(rdb *redis.Client) error {
	if rdb == nil {
		return fmt.Errorf("koneksi Redis tidak tersedia")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}

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

	if serverInfo, err := rdb.Info(ctx, "server").Result(); err == nil {
		result["server"] = serverInfo
	}

	return result, nil
}
