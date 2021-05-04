package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/vdobler/chart"

	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/pronote/api"
)

func TimetableChartCommand() lib.Command {
	return lib.Command{
		Name:        "timetable_chart",
		Description: "Cette commande permet de tracer un diagramme en camembert (ou en quartiers) afin de visualiser le volume horaire des matières dans l'emploi du temps.",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, _ []string, _ map[string]interface{}) error {
			response, err := api.GetTimetable(now.BeginningOfWeek(), now.EndOfWeek())
			if err != nil {
				return lib.Error(bot, update, "Erreur serveur : impossible d'effectuer la requête vers PRONOTE.")
			}

			if len(response.Timetable) == 0 {
				msg := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun cours n'est prévu pour le moment.")
				_, err = bot.Send(msg)
				return err
			}

			data := make(map[string]float64)
			var totalDuration float64
			for _, lesson := range response.Timetable {
				duration := float64(lesson.To - lesson.From)

				_, exists := data[lesson.Subject]
				if !exists {
					data[lesson.Subject] = duration
				} else {
					data[lesson.Subject] += duration
				}
				totalDuration += duration
			}

			var subjects []string
			var durations []float64
			for subject, duration := range data {
				subjects = append(subjects, subject)
				durations = append(durations, (duration/totalDuration)*100)
			}

			pie := chart.PieChart{Title: "Distribution des cours de la semaine"}
			pie.AddDataPair("Cours", subjects, durations)
			pie.Inner = 0.7
			pie.FmtVal = chart.PercentValue

			file := lib.Plot(&pie, "timetable_chart")
			photo := telegram.NewPhotoUpload(update.Message.Chat.ID, file)
			_, err = bot.Send(photo)
			return err
		},
	}
}
