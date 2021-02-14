package math

import (
	"github.com/Knetic/govaluate"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vdobler/chart"
	"image/color"
	"strings"

	"github.com/theovidal/bacbot/lib"
)

func PlotCommand() lib.Command {
	return lib.Command{
		Flags: map[string]interface{}{
			"x_min":   -10.0,
			"x_max":   10.0,
			"x_scale": 1.0,
			"y_min":   -10.0,
			"y_max":   10.0,
			"y_scale": 1.0,

			"color":      "red",
			"line_width": 1,
		},
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error {
			lineColor, exists := lib.Colors[flags["color"].(string)]
			if !exists {
				return lib.Error(bot, update, "La couleur spécifiée n'existe pas. Vérifiez la liste des couleurs disponibles sur la page d'aide de la commande.")
			}

			function := strings.Join(args, " ")

			plot := chart.ScatterChart{Title: "f(x) = " + function}
			plot.XRange.MinMode.Fixed, plot.XRange.MaxMode.Fixed = true, true
			plot.XRange.MinMode.Value, plot.XRange.MaxMode.Value = flags["x_min"].(float64), flags["x_max"].(float64)

			plot.YRange.MinMode.Fixed, plot.YRange.MaxMode.Fixed = true, true
			plot.YRange.MinMode.Value, plot.YRange.MaxMode.Value = flags["y_min"].(float64), flags["y_max"].(float64)

			plot.XRange.TicSetting.Delta = flags["x_scale"].(float64)
			plot.YRange.TicSetting.Delta = flags["y_scale"].(float64)

			plot.XRange.TicSetting.Mirror = -1
			plot.YRange.TicSetting.Mirror = -1

			style := chart.Style{
				Symbol:    'o',
				LineColor: lineColor,
				FillColor: color.NRGBA{0xff, 0x80, 0x80, 0xff},
				LineStyle: chart.SolidLine,
				LineWidth: flags["line_width"].(int),
			}

			_, err := govaluate.NewEvaluableExpressionWithFunctions(function, functions)
			if err != nil {
				return lib.Error(bot, update, "La fonction entrée est invalide : `"+err.Error()+"`.")
			}

			config := telegram.NewMessage(update.Message.Chat.ID, "_Génération du graphique en cours..._")
			config.ParseMode = "Markdown"
			waiter, _ := bot.Send(config)

			plot.AddFunc("C_f", func(x float64) float64 {
				expression, _ := govaluate.NewEvaluableExpressionWithFunctions(function, functions)
				y, _ := expression.Evaluate(map[string]interface{}{"x": x})
				return y.(float64)
			}, chart.PlotStyleLines, style)

			file := lib.Plot(&plot, "function_plot")
			photo := telegram.NewPhotoUpload(update.Message.Chat.ID, file)
			_, err = bot.Send(photo)
			bot.DeleteMessage(telegram.DeleteMessageConfig{
				ChatID:    waiter.Chat.ID,
				MessageID: waiter.MessageID,
			})
			return err
		},
	}
}
