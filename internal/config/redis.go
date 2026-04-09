package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(uri string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: uri,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	if RedisClient == nil {
		log.Fatal("Redis client is nil.")
	}
	log.Println("Redis connected.")
}
