package config

import (
	"context"
	"log"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Имя контейнера из docker-compose
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
	log.Println("✅ Redis подключен")
}

func RedisClient() *redis.Client {
	return redisClient
}
