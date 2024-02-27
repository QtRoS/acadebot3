package searchengine

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis"
	"github.com/qtros/acadebot3/internal/utils"
)

const searchTTLMinutes = 60

func init() {
	utils.CommonClient.Timeout = utils.CommonClient.Timeout + 2*time.Second
}

// RudraSearch for courses in Rudra.
func RudraSearch(query string, limit int) string {
	redisKey := fmt.Sprintf("query:%x", md5.Sum([]byte(query)))

	value, err := utils.RedisClient.Get(redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			slog.Info("Redis miss: ", slog.Any("query", query))
		} else {
			slog.Error("Redis error:", slog.Any("err", err))
		}

		newValue := Search(query, limit)
		if newValue != "" {
			utils.RedisClient.Set(redisKey, newValue, time.Minute*searchTTLMinutes)
		}
		value = newValue
	} else {
		slog.Info("Rudra Redis hit", slog.Any("redisKey", redisKey))
	}

	return value
}
