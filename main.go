package main

import (
	"log"
	"os"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/parcoursup"
	"github.com/theovidal/bacbot/pronote"
)

func main() {
	lib.LoadEnv(".env")
	lib.OpenCache()

	bot, err := telegram.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Println(lib.Green.Sprintf("✅ Authorized on account %s", bot.Self.UserName))

	go func() {
		for range time.Tick(time.Minute * 10) {
			err := pronote.TimetableTicker(bot)
			if err != nil {
				log.Println(lib.Red.Sprintf("‼ Error handling timetable ticker: %s", err))
			}
		}
	}()

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {
			parcoursup.HandleRequest(bot, &update)
		} else if update.Message.IsCommand() {
			if update.Message.From.UserName != os.Getenv("TELEGRAM_USER") {
				continue
			}
			err := HandleCommand(bot, update, false)
			if err != nil {
				log.Println(lib.Red.Sprintf("‼ Error handling an event: %s", err))
			}
		}
	}
}
