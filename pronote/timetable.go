package pronote

import (
	"fmt"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

func TimetableCommand() lib.Command {
	return lib.Command{
		Name:        "timetable",
		Description: "Cette commande permet d'obtenir l'emploi du temps complet sur les 7 prochains jours, avec leur statut à jour et le mode présentiel/distanciel.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			response, err := api.GetTimetable(time.Now(), time.Now().Add(time.Hour*24*6))
			if err != nil {
				lib.LogError(err.Error())
				return lib.Error(bot, chatID, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Timetable) == 0 {
				msg := telegram.NewMessage(chatID, "🍃 Aucun cours n'est prévu pour le moment.")
				_, err = bot.Send(msg)
				return
			}

			days := make(map[string]string)
			for _, lesson := range response.Timetable {
				day := time.Unix(int64(lesson.From/1000), 0).Format("02/01")

				days[day] = days[day] + lesson.String()
			}

			var content string
			for day, lessons := range days {
				content += fmt.Sprintf("*―――――― %s ――――――*\n%s\n", day, lessons)
			}

			msg := telegram.NewMessage(chatID, content)
			msg.ParseMode = "MarkdownV2"
			_, err = bot.Send(msg)
			return
		},
	}
}
