package acadebot

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/qtros/acadebot3/internal/models"
	"github.com/qtros/acadebot3/internal/searchengine"
	"github.com/qtros/acadebot3/internal/utils"
)

const (
	searchTTLMinutes  = 60
	contextTTLMinutes = 15
)

func getCourses(query string) []models.CourseInfo {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return nil
	}

	jsonStr := searchengine.RudraSearch(query, perSrcLimit)
	//slog.Debug(jsonStr)

	var courses []models.CourseInfo
	if err := json.Unmarshal([]byte(jsonStr), &courses); err != nil {
		slog.Error("Bad JSON:", slog.Any("err", err))
		return nil
	}

	// for _, c := range courses {
	// 	slog.Debug(c.Link)
	// }

	return courses
}

func courseInfoToInlineQueryResult(c models.CourseInfo) tgbotapi.InlineQueryResultArticle {
	id := fmt.Sprintf("%x", md5.Sum([]byte(c.Link)))
	article := tgbotapi.NewInlineQueryResultArticle(id, c.Name, c.String())
	article.URL = c.Link
	article.ThumbURL = c.Art
	return article
}

func createKeyboard(context *models.UserContext) tgbotapi.InlineKeyboardMarkup {
	status := fmt.Sprintf("%d of %d", context.Position+1, context.Count)
	btbl := tgbotapi.NewInlineKeyboardButtonData(status, "0")

	bm := tgbotapi.NewInlineKeyboardButtonData("◀ Previous", "-1")
	bp := tgbotapi.NewInlineKeyboardButtonData("Next   ▶", "+1")
	bfm := tgbotapi.NewInlineKeyboardButtonData("⏪", "-5")
	bfp := tgbotapi.NewInlineKeyboardButtonData("⏩", "+5")

	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(bm, bp),
		tgbotapi.NewInlineKeyboardRow(bfm, btbl, bfp))
}

// saveContext to Redis.
func saveContext(chatID int64, context *models.UserContext) {
	redisKey := fmt.Sprintf("usercontext:%d", chatID)
	value, err := json.Marshal(context)
	if err != nil {
		slog.Error("marshal error", slog.Any("err", err))
		return
	}

	utils.RedisClient.Set(redisKey, value, time.Minute*contextTTLMinutes)
}

// restoreContext from Redis.
func restoreContext(chatID int64) *models.UserContext {
	redisKey := fmt.Sprintf("usercontext:%d", chatID)

	value, err := utils.RedisClient.Get(redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			slog.Warn("No context for:", slog.Any("redisKey", redisKey))
		} else {
			slog.Error("Redis error:", slog.Any("err", err))
		}
		return nil
	}

	var context models.UserContext
	err = json.Unmarshal([]byte(value), &context)
	if err != nil {
		slog.Error("marshal error", slog.Any("err", err))
	}
	return &context
}
