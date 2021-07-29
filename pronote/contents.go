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
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, _ []string, _ map[string]interface{}) (err error) {
			response, err := api.GetContents()
			if err != nil {
				lib.LogError(err.Error())
				return lib.Error(bot, chatID, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Contents) == 0 {
				msg := telegram.NewMessage(chatID, "🍃 Aucun contenu de cours n'a été saisi pour le moment.")
				_, err = bot.Send(msg)
				return
			}

			var output []string
			for _, content := range response.Contents.Reverse() {
				output = append(output, content.String())
			}

			msg := telegram.NewMessage(chatID, strings.Join(output, "―――――――――\n"))
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			_, err = bot.Send(msg)
			return
		},
	}
}
