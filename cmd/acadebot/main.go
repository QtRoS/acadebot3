package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/qtros/acadebot3/internal/acadebot"
)

const (
	envAPIKey        = "ENV_API_KEY"
	perSrcLimit      = 10
	noCoursesFound   = "Sorry, no similar course found."
	dummyPlaceholder = "Just a moment..."
	noContextFound   = "Sorry, can't navigate through results. Try to search again!"
	greeting         = `Hello, %s!
I can help you with finding online courses (MOOCs).
Type course name or keyword and I will find something for you! 
	https://github.com/QtRoS/acadebot3`
)

var bot *tgbotapi.BotAPI

func init() {
	token := os.Getenv(envAPIKey)
	slog.Info("Token:", slog.Any("token", token)) // TODO BUG

	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	slog.Info(fmt.Sprintf("Authorized on account %s", b.Self.UserName))
	// bot.Debug = true

	bot = b
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.InlineQuery != nil {
			acadebot.HandleInlineQuery(bot, update.InlineQuery)
		} else if update.CallbackQuery != nil {
			acadebot.HandleCallbackQuery(bot, update.CallbackQuery)
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				acadebot.HandleCommand(bot, update.Message)
			} else {
				acadebot.HandleMessage(bot, update.Message)
			}
		}
	}
}
