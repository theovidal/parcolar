package main

import (
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcoursupbot/pronote"
)

var commandsList = map[string]func(bot *telegram.BotAPI, update telegram.Update, args []string) error{
	"homework":  pronote.HomeworkCommand,
	"timetable": pronote.TimetableCommand,
}

func HandleCommand(bot *telegram.BotAPI, update telegram.Update, isCallback bool) error {
	var commandName string
	var args []string

	if isCallback {
		parts := strings.Split(strings.TrimPrefix(update.CallbackQuery.Data, "/"), " ")
		commandName = parts[0]
		args = parts[1:]
	} else {
		commandName = update.Message.Command()
		if update.Message.CommandArguments() != "" {
			args = strings.Split(update.Message.CommandArguments(), " ")
		}
	}

	command, exists := commandsList[commandName]
	if !exists {
		_, err := bot.Send(telegram.NewMessage(update.Message.Chat.ID, "❓ Oups, il semble que cette commande soit inconnue!"))
		return err
	}

	err := command(bot, update, args)
	return err
}
