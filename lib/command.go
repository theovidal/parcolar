package lib

import telegram "github.com/go-telegram-bot-api/telegram-bot-api"

type Command struct {
	Description string
	Flags       map[string]interface{}
	Execute     func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error
}
