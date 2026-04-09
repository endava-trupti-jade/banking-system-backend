package redisrepo

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/config"
	"context"
	"log"
	"time"
)

type CacheRepository struct{}

func NewCacheRepository() *CacheRepository {
	return &CacheRepository{}
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	log.Println("CacheRepository Get() started")
	if config.RedisClient == nil {
		return "", constants.ErrRedisUninitialized
	}

	log.Println("CacheRepository Get() end")
	return config.RedisClient.Get(ctx, key).Result()
}

func (r *CacheRepository) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	log.Println("CacheRepository Set() started")
	if ttl == 0 {
		ttl = time.Hour * 24 // default TTL
	}

	if config.RedisClient == nil {
		return constants.ErrRedisUninitialized
	}

	log.Println("CacheRepository Set() end")
	return config.RedisClient.Set(ctx, key, value, ttl).Err()
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	log.Println("CacheRepository Delete() started")
	if config.RedisClient == nil {
		return constants.ErrRedisUninitialized
	}

	log.Println("CacheRepository Delete() end")
	return config.RedisClient.Del(ctx, key).Err()
}
