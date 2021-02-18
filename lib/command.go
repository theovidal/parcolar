package lib

import telegram "github.com/go-telegram-bot-api/telegram-bot-api"

// Command describes a Telegram commands with all its information
type Command struct {
	// A complete description of the command to show in the help message
	Description string
	// The list of flags that can be passed by the user
	Flags map[string]Flag
	// The command actions
	Execute func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error
}

type Flag struct {
	Description string
	Value       interface{}
}
