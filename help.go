package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
)

func HelpCommand() lib.Command {
	return lib.Command{
		Name:        "help",
		Description: "Obtenir de l'aide sur le bot et ses commandes. Envoyez `/help` suivi du nom de la commande recherchée pour en apprendre davantage.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				message := telegram.NewMessage(chatID, "*―――――― 🍀 Parcolar ――――――*\n\nCe bot Telegram propose de nombreux outils conçus dans un but éducatif mais qui peut également être utilisé dans n'importe quel contexte. Commencez par taper `/` et découvrez la liste des commandes disponibles. Pour obtenir davantage d'informations sur l'une d'entre elles, tapez `/help` suivi du nom de la commande.")
				message.ParseMode = "Markdown"
				_, err = bot.Send(message)
				return
			}

			command, found := commandsList[args[0]]
			if !found {
				return lib.Error(bot, chatID, "La commande recherchée est inconnue. Vérifiez que vous n'ayez pas fait de fautes de frappe.")
			}

			msg := telegram.NewMessage(chatID, command.Help())
			msg.ParseMode = "Markdown"
			_, err = bot.Send(msg)
			return
		},
	}
}
