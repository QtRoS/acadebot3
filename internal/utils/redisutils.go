package utils

import (
	"log/slog"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

const (
	envRedisAddress   = "ENV_REDIS_ADDRESS"
	envRedisPass      = "ENV_REDIS_PASS"
	searchTTLMinutes  = 60
	contextTTLMinutes = 15
)

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv(envRedisAddress) + ":6379",
		Password: os.Getenv(envRedisPass), // no password set
		DB:       0,                       // use default DB
	})

	pong, err := RedisClient.Ping().Result()
	slog.Info("Redis ping:", slog.Any("pong", pong), slog.Any("err", err))
}
