package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/vdobler/chart"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

func TimetableChartCommand() lib.Command {
	return lib.Command{
		Name:        "timetable_chart",
		Description: "Cette commande permet de tracer un diagramme en camembert (ou en quartiers) afin de visualiser le volume horaire des matières dans l'emploi du temps.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, _ []string, _ map[string]interface{}) (err error) {
			response, err := api.GetTimetable(now.BeginningOfWeek(), now.EndOfWeek())
			if err != nil {
				return lib.Error(bot, update, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Timetable) == 0 {
				message := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun cours n'est prévu pour le moment.")
				_, err = bot.Send(message)
				return
			}

			lessons := make(map[string]float64)
			var totalDuration float64
			for _, lesson := range response.Timetable {
				duration := float64(lesson.To - lesson.From)

				_, exists := lessons[lesson.Subject]
				if !exists {
					lessons[lesson.Subject] = duration
				} else {
					lessons[lesson.Subject] += duration
				}
				totalDuration += duration
			}

			var subjects []string
			var durations []float64
			for subject, duration := range lessons {
				subjects = append(subjects, subject)
				durations = append(durations, (duration/totalDuration)*100)
			}

			pieChart := chart.PieChart{Title: "Distribution des cours de la semaine"}
			pieChart.AddDataPair("Cours", subjects, durations)
			pieChart.Inner = 0.7
			pieChart.FmtVal = chart.PercentValue

			file := lib.Plot(&pieChart, "timetable_chart")
			photoUpload := telegram.NewPhotoUpload(update.Message.Chat.ID, file)
			_, err = bot.Send(photoUpload)
			return
		},
	}
}
