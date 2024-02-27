package searchengine

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/qtros/acadebot3/internal/models"
)

const (
	envRedisAddress = "ENV_REDIS_ADDRESS"
	envRedisPass    = "ENV_REDIS_PASS"
)

type cachingAdapter struct {
	sourceAdapter SourceAdapter
	client        *redis.Client
	ttl           time.Duration
}

func newCachingAdapter(adapter SourceAdapter, ttl time.Duration) *cachingAdapter {
	ca := cachingAdapter{}
	ca.sourceAdapter = adapter
	ca.ttl = ttl

	ca.client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv(envRedisAddress) + ":6379",
		Password: os.Getenv(envRedisPass), // no password set
		DB:       0,                       // use default DB
	})

	pong, err := ca.client.Ping().Result()
	slog.Info("Caching client", ca.sourceAdapter.Name(), "redis ping: ", pong, err)

	return &ca
}

func (me *cachingAdapter) Name() string {
	return me.sourceAdapter.Name() + " (Cached)"
}

func (me *cachingAdapter) Get(query string, limit int) []models.CourseInfo {
	redisKey := fmt.Sprintf("cachingadapter:%s", me.sourceAdapter.Name())

	value, err := me.client.Get(redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			slog.Info("Redis miss:", "redisKey", redisKey)
		} else {
			slog.Error("Redis error:", slog.Any("err", err))
		}

		rawData := me.sourceAdapter.Get(query, limit)
		rawDataAsJSON, err := json.Marshal(rawData)
		if err != nil {
			slog.Error("marshal error", slog.Any("err", err))
		} else if rawData != nil {
			me.client.Set(redisKey, rawDataAsJSON, me.ttl)
		}

		return rawData
	}

	slog.Info("Redis HIT:", "redisKey", redisKey)
	var result []models.CourseInfo
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		slog.Error("marshal error", slog.Any("err", err))
	}

	return result
}
