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
		Description: "Obtenir les contenus des cours des 5 derniers jours.",
		Flags: map[string]lib.Flag{
			"days": {"Nombre de jours en arrière (sans compter aujourd'hui)", 4, nil},
		},
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, _ []string, flags map[string]interface{}) (err error) {
			response, err := api.GetContents(bot.Cache, flags["days"].(int))
			if err != nil {
				lib.LogError(err.Error())
				return bot.Error(chatID, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
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
