package main

import (
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/theovidal/bacbot/lib"
)

func HelpCommand() lib.Command {
	return lib.Command{
		Description: "Obtenez l'aide détaillée sur les commandes du bot.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error {
			if len(args) == 0 {
				return lib.Error(bot, update, "Merci d'ajouter le nom de la commande recherchée au message.")
			}

			command, found := commandsList[args[0]]
			if !found {
				return lib.Error(bot, update, "La commande recherchée est inconnue.")
			}

			content := fmt.Sprintf("*―――――― Aide de la commande %s ――――――*\n%s", args[0], command.Description)

			if len(command.Flags) > 0 {
				content += "\n\nListe des flags disponibles sur cette commande :\n"
				for name, flag := range command.Flags {
					content += fmt.Sprintf("• `%s` : %s _(par défaut : %v)_\n", name, flag.Description, flag.Value)
				}
				content += "Les flags sont à ajouter en début de commande sous la forme `nom=valeur`. Veillez à respecter le type de chacun (nombre entier, réel...)"
			}

			msg := telegram.NewMessage(update.Message.Chat.ID, content)
			msg.ParseMode = "Markdown"
			_, err := bot.Send(msg)
			return err
		},
	}
}
