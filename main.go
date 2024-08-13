package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"golang/echo-bot/utils"
)

var Keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
	),
)

var (
	todayMin  = 58
	todayHour = 20
)

const (
	minHour = 7
	maxHour = 23
	minMin  = 0
	maxMin  = 59
)

func commandHandler(bot *tgbotapi.BotAPI, command *tgbotapi.Message) {
	user := command.From
	text := command.Text
	switch text {
	case "/start":
		msg := tgbotapi.NewMessage(user.ID, fmt.Sprintf("Hello %s", user))
		msg.ReplyMarkup = Keyboard
		bot.Send(msg)
		go CheckTime(bot, user)
	}
}

func Predict(bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	randPred := utils.Random(0, len(utils.Predictions)-1)
	msg := tgbotapi.NewMessage(user.ID, utils.Predictions[randPred])
	msg.ReplyMarkup = Keyboard
	bot.Send(msg)
}

func CheckTime(bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	for {
		log.Print(time.Now())
		if time.Now().Hour() == todayHour && time.Now().Minute() == todayMin {
			Predict(bot, user)
			todayMin = utils.Random(minMin, maxMin)
			todayHour = utils.Random(minHour, maxHour)
		}
		time.Sleep(59 * time.Second)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
	token, _ := os.LookupEnv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			commandHandler(bot, update.Message)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
