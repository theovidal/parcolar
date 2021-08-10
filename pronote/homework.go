package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

func HomeworkCommand() lib.Command {
	return lib.Command{
		Name:        "homework",
		Description: "Obtenir tous les devoirs saisis sur PRONOTE pour les prochains jours.",
		Flags: map[string]lib.Flag{
			"days": {"Nombre de jours à obtenir (sans compter aujourd'hui)", 14, nil},
		},
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, _ []string, flags map[string]interface{}) (err error) {
			response, err := api.GetHomework(bot.Cache, flags["days"].(int))
			if err != nil {
				lib.LogError(err.Error())
				return bot.Error(chatID, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Homeworks) == 0 {
				msg := telegram.NewMessage(chatID, "🍃 Aucun devoir n'a été rédigé pour le moment.")
				_, err = bot.Send(msg)
				return
			}

			content := ""
			for _, homework := range response.Homeworks {
				content += homework.String()
			}

			msg := telegram.NewMessage(chatID, content)
			msg.ParseMode = "MarkdownV2"
			msg.DisableWebPagePreview = true
			_, err = bot.Send(msg)
			return
		},
	}
}
