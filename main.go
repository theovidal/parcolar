package main

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcoursupbot/lib"
	"github.com/theovidal/parcoursupbot/parcoursup"
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

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {
			parcoursup.HandleRequest(bot, &update)
		} else if update.Message.IsCommand() {
			err := HandleCommand(bot, update, false)
			if err != nil {
				log.Println(lib.Red.Sprintf("‼ Error handling an event: %s", err))
			}
		}
	}
}
