package acadebot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/qtros/acadebot3/internal/models"
)

const (
	perSrcLimit      = 10
	noCoursesFound   = "Sorry, no similar course found."
	dummyPlaceholder = "Just a moment..."
	noContextFound   = "Sorry, can't navigate through results. Try to search again!"
	greeting         = `Hello, %s!
I can help you with finding online courses (MOOCs).
Type course name or keyword and I will find something for you! 
	https://github.com/QtRoS/acadebot3`
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	slog.Info("Message [%s] %s", message.From.UserName, message.Text)

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, dummyPlaceholder))

	query := strings.TrimSpace(message.Text)
	courses := getCourses(query)
	if courses == nil || len(courses) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, noCoursesFound)
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	} else {
		context := models.UserContext{Query: query, Position: 0, Count: len(courses)}
		saveContext(message.Chat.ID, &context)

		courseInfo := courses[context.Position]

		msg := tgbotapi.NewMessage(message.Chat.ID, courseInfo.String())
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = createKeyboard(&context)
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

func HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	var answer string
	switch command := message.Command(); command {
	case "start":
		answer = fmt.Sprintf(greeting, message.From.UserName)
	default:
		answer = fmt.Sprintf("Unknown command: %s", command)
	}

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, answer))
}

func HandleInlineQuery(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery) {
	slog.Info("Inline [%s] %s", inlineQuery.From.UserName, inlineQuery.Query)
	courses := getCourses(inlineQuery.Query)
	if len(courses) == 0 {
		return
	}

	var articles = make([]interface{}, len(courses))
	for i, c := range courses {
		article := courseInfoToInlineQueryResult(c)
		articles[i] = article
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}

	slog.Debug(fmt.Sprintf("Article count: %d", len(articles)))
	// if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
	if _, err := bot.Request(inlineConf); err != nil {
		slog.Error("can't send answerInlineQuery", slog.Any("err", err))
	}
}

func HandleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	// Dummy answer to stop spinners in UI.
	// bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: callbackQuery.ID})
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := bot.Request(callback); err != nil {
		slog.Warn("can't send answerCallbackQuery", err)
	}

	// Check if there context for that user.
	context := restoreContext(callbackQuery.Message.Chat.ID)
	if context == nil {
		bot.Send(tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, noContextFound))
		return
	}

	// Calculate delta.
	delta, _ := strconv.Atoi(callbackQuery.Data)
	context.Position = max(context.Position+delta, 0)
	context.Position = min(context.Position, context.Count-1)
	saveContext(callbackQuery.Message.Chat.ID, context)

	// Get last results.
	courses := getCourses(context.Query)
	if len(courses) == 0 {
		slog.Warn("getCourses returned no courses")
		return
	}

	courseInfo := courses[min(len(courses)-1, context.Position)]

	// Answer in TG.
	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID,
		callbackQuery.Message.MessageID, courseInfo.String())
	keyboard := createKeyboard(context)
	msg.BaseEdit.ReplyMarkup = &keyboard
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
