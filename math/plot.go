package math

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vdobler/chart"

	"github.com/theovidal/bacbot/lib"
)

func PlotCommand(bot *telegram.BotAPI, update telegram.Update, args []string) error {
	chartParams := map[string]float64{
		"x_min":   -10,
		"x_max":   10,
		"x_scale": 1,
		"y_min":   -10,
		"y_max":   10,
		"y_scale": 1,
	}
	var function string
	for index, arg := range args {
		if strings.Contains(arg, "=") {
			flag := strings.Split(arg, "=")
			name := flag[0]

			value, err := strconv.ParseFloat(flag[1], 64)
			if err != nil {
				msg := telegram.NewMessage(update.Message.Chat.ID, "❌ Merci d'indiquer des nombres décimaux en tant que paramètres du graphique.")
				_, err := bot.Send(msg)
				return err
			}
			chartParams[name] = value
		} else {
			function = strings.Join(args[index:], " ")
			break
		}
	}

	plot := chart.ScatterChart{Title: "f(x) = " + function}
	plot.XRange.MinMode.Fixed, plot.XRange.MaxMode.Fixed = true, true
	plot.XRange.MinMode.Value, plot.XRange.MaxMode.Value = chartParams["x_min"], chartParams["x_max"]

	plot.YRange.MinMode.Fixed, plot.YRange.MaxMode.Fixed = true, true
	plot.YRange.MinMode.Value, plot.YRange.MaxMode.Value = chartParams["y_min"], chartParams["y_max"]

	plot.XRange.TicSetting.Delta = chartParams["x_scale"]
	plot.YRange.TicSetting.Delta = chartParams["y_scale"]

	plot.XRange.TicSetting.Mirror = -1
	plot.YRange.TicSetting.Mirror = -1

	style := chart.Style{
		Symbol:    'o',
		LineColor: lib.Colors["red"],
		FillColor: color.NRGBA{0xff, 0x80, 0x80, 0xff},
		LineStyle: chart.SolidLine,
		LineWidth: 2,
	}

	_, err := govaluate.NewEvaluableExpressionWithFunctions(function, functions)
	if err != nil {
		msg := telegram.NewMessage(update.Message.Chat.ID, "❌ La fonction entrée est invalide : `"+err.Error()+"`.")
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		return err
	}

	config := telegram.NewMessage(update.Message.Chat.ID, "_Génération du graphique en cours..._")
	config.ParseMode = "Markdown"
	msg, _ := bot.Send(config)

	plot.AddFunc("C_f", func(x float64) float64 {
		expression, _ := govaluate.NewEvaluableExpressionWithFunctions(function, functions)

		y, _ := expression.Evaluate(map[string]interface{}{"x": x})
		return y.(float64)
	}, chart.PlotStyleLines, style)

	file := lib.Plot(&plot, "function_plot")
	photo := telegram.NewPhotoUpload(update.Message.Chat.ID, file)
	_, err = bot.Send(photo)
	bot.DeleteMessage(telegram.DeleteMessageConfig{
		ChatID:    msg.Chat.ID,
		MessageID: msg.MessageID,
	})
	return err
}
