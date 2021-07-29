package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
)

func HelpCommand() lib.Command {
	return lib.Command{
		Name:        "help",
		Description: "Obtenez l'aide détaillée sur les commandes du bot.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				return lib.Error(bot, chatID, "Merci d'ajouter le nom de la commande recherchée au message.")
			}

			command, found := commandsList[args[0]]
			if !found {
				return lib.Error(bot, chatID, "La commande recherchée est inconnue.")
			}

			msg := telegram.NewMessage(chatID, command.Help())
			msg.ParseMode = "Markdown"
			_, err = bot.Send(msg)
			return
		},
	}
}
