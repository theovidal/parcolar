package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/theovidal/parcolar/math/src"
	"github.com/vdobler/chart"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

func TimetableChartCommand() lib.Command {
	return lib.Command{
		Name:        "timetable_chart",
		Description: "Tracer un diagramme en camembert (ou en quartiers) pour visualiser le volume horaire des matières dans l'emploi du temps de la semaine en cours.",
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, _ []string, _ map[string]interface{}) (err error) {
			response, err := api.GetTimetable(bot.Cache, now.BeginningOfWeek(), now.EndOfWeek())
			if err != nil {
				return bot.Error(chatID, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Timetable) == 0 {
				message := telegram.NewMessage(chatID, "🍃 Aucun cours n'est prévu pour le moment.")
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

			file := src.Plot(&pieChart, "timetable_chart")
			photoUpload := telegram.NewPhotoUpload(chatID, file)
			_, err = bot.Send(photoUpload)
			return
		},
	}
}
