package math

import (
	"fmt"
	"image/color"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vdobler/chart"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/math/src"
)

func PlotCommand() lib.Command {
	return lib.Command{
		Name:        "plot",
		Description: fmt.Sprintf("Tracer des graphiques riches et personnalisés. Vous pouvez tracer plusieurs fonctions en séparant leurs expressions par une esperluette `&`.\n%s\n\n%s", src.DataDocumentation, src.CalcDisclaimer),
		Flags: map[string]lib.Flag{
			"x_min":   {"Valeur minimale de `x` à afficher", -10.0, nil},
			"x_max":   {"Valeur maximale de `x` à afficher", 10.0, nil},
			"x_scale": {"Pas pour l'abscisse", 1.0, nil},
			"y_min":   {"Valeur minimale de `y` à afficher", -10.0, nil},
			"y_max":   {"Valeur maximale de `y` à afficher", 10.0, nil},
			"y_scale": {"Pas pour l'ordonnée", 1.0, nil},

			// TODO: color choice with multiple functions
			// "color":      {"Couleur de la courbe : `red`, `pink`, `purple`, `indigo`, `blue`, `light_blue`, `cyan`, `teal`, `green`, `light_green`, `lime`, `yellow`, `amber`, `orange`, `brown`.", "red"},
			"line_width": {"Épaisseur de la courbe (en pixels)", 1, nil},

			"grid": {"Afficher la grille sur le graphique (0 ou 1)", true, nil},
		},
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) error {
			if len(args) == 0 {
				return bot.Help(chatID, "plot")
			}

			grid := flags["grid"].(bool)

			raw := strings.Join(args, " ")
			functions := strings.Split(raw, "&")

			plot := chart.ScatterChart{}
			plot.XRange.MinMode.Fixed, plot.XRange.MaxMode.Fixed = true, true
			plot.XRange.MinMode.Value, plot.XRange.MaxMode.Value = flags["x_min"].(float64), flags["x_max"].(float64)

			plot.YRange.MinMode.Fixed, plot.YRange.MaxMode.Fixed = true, true
			plot.YRange.MinMode.Value, plot.YRange.MaxMode.Value = flags["y_min"].(float64), flags["y_max"].(float64)

			if grid {
				plot.XRange.TicSetting.Grid = chart.GridLines
				plot.YRange.TicSetting.Grid = chart.GridLines
			}

			plot.XRange.TicSetting.Delta = flags["x_scale"].(float64)
			plot.YRange.TicSetting.Delta = flags["y_scale"].(float64)

			plot.XRange.TicSetting.Mirror = -1
			plot.YRange.TicSetting.Mirror = -1

			style := chart.Style{
				Symbol:    'o',
				FillColor: color.NRGBA{R: 0xff, G: 0x80, B: 0x80, A: 0xff},
				LineStyle: chart.SolidLine,
				LineWidth: flags["line_width"].(int),
			}

			for _, function := range functions {
				if err := src.CheckExpression(function); err != nil {
					return bot.Error(chatID, err.Error())
				}
				if _, err := src.Evaluate(function, 1); err != nil {
					return bot.Error(chatID, err.Error())
				}
			}

			config := telegram.NewMessage(chatID, "_Génération du graphique en cours..._")
			config.ParseMode = "Markdown"
			waiter, _ := bot.Send(config)

			colorNumber := 0
			for _, function := range functions {
				current := strings.TrimSpace(function)
				if colorNumber == len(src.PlotColors) {
					colorNumber = 0
				}
				style.LineColor = src.PlotColors[colorNumber]
				colorNumber++

				plot.AddFunc(current, func(x float64) float64 {
					value, _ := src.Evaluate(current, x)
					return value
				}, chart.PlotStyleLines, style)
			}

			file := src.Plot(&plot, "function_plot")
			photo := telegram.NewPhotoUpload(chatID, file)
			_, err := bot.Send(photo)
			if err != nil {
				return err
			}
			_, err = bot.DeleteMessage(telegram.DeleteMessageConfig{
				ChatID:    waiter.Chat.ID,
				MessageID: waiter.MessageID,
			})
			return err
		},
	}
}
