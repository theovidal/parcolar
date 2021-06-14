package pronote

import (
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

func ContentsCommand() lib.Command {
	return lib.Command{
		Name:        "contents",
		Description: "Cette commande permet d'obtenir les contenus des cours des 5 derniers jours.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, _ []string, _ map[string]interface{}) error {
			response, err := api.GetContents()
			if err != nil {
				return lib.Error(bot, update, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Contents) == 0 {
				msg := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun contenu de cours n'a été saisi pour le moment.")
				_, err = bot.Send(msg)
				return err
			}

			var output []string
			for _, content := range response.Contents.Reverse() {
				output = append(output, content.String())
			}

			msg := telegram.NewMessage(update.Message.Chat.ID, strings.Join(output, "―――――――――\n"))
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			_, err = bot.Send(msg)

			return err
		},
	}
}
