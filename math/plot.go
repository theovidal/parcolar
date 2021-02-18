package math

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/Knetic/govaluate"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vdobler/chart"

	"github.com/theovidal/bacbot/lib"
)

func PlotCommand() lib.Command {
	var functionsDescription string
	for name := range functions {
		functionsDescription += fmt.Sprintf("`%s`, ", name)
	}
	functionsDescription = strings.TrimSuffix(functionsDescription, ", ")

	var constantsDescription string
	for name := range constants {
		constantsDescription += fmt.Sprintf("`%s`, ", name)
	}
	constantsDescription = strings.TrimSuffix(constantsDescription, ", ")

	return lib.Command{
		Description: fmt.Sprintf("Tracer des graphiques riches et complets.\n📈 Les fonctions disponibles sont : %s.\nπ Les constantes disponibles sont: %s.\n\n⚠ *Tous les signes multiplier* sont obligatoires (ex: 3x => 3 \\* x) et les *puissances* sont représentées par une *double-étoile* (\\*\\*).", functionsDescription, constantsDescription),
		Flags: map[string]lib.Flag{
			"x_min":   {"Valeur minimale de `x` à afficher", -10.0},
			"x_max":   {"Valeur maximale de `x` à afficher", 10.0},
			"x_scale": {"Pas pour l'abscisse", 1.0},
			"y_min":   {"Valeur minimale de `y` à afficher", -10.0},
			"y_max":   {"Valeur maximale de `y` à afficher", 10.0},
			"y_scale": {"Pas pour l'ordonnée", 1.0},

			"color":      {"Couleur de la courbe : `red`, `pink`, `purple`, `indigo`, `blue`, `light_blue`, `cyan`, `teal`, `green`, `light_green`, `lime`, `yellow`, `amber`, `orange`, `brown`.", "red"},
			"line_width": {"Épaisseur de la courbe (en pixels)", 1},

			"grid": {"Afficher la grille sur le graphique (0 ou 1)", 1},
		},
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error {
			lineColor, exists := lib.Colors[flags["color"].(string)]
			if !exists {
				return lib.Error(bot, update, "La couleur spécifiée n'existe pas. Vérifiez la liste des couleurs disponibles sur la page d'aide de la commande.")
			}
			grid := flags["grid"].(int)

			function := strings.Join(args, " ")

			plot := chart.ScatterChart{Title: "f(x) = " + function}
			plot.XRange.MinMode.Fixed, plot.XRange.MaxMode.Fixed = true, true
			plot.XRange.MinMode.Value, plot.XRange.MaxMode.Value = flags["x_min"].(float64), flags["x_max"].(float64)
			if grid == 1 {
				plot.XRange.TicSetting.Grid = chart.GridLines
			}

			plot.YRange.MinMode.Fixed, plot.YRange.MaxMode.Fixed = true, true
			plot.YRange.MinMode.Value, plot.YRange.MaxMode.Value = flags["y_min"].(float64), flags["y_max"].(float64)
			if grid == 1 {
				plot.YRange.TicSetting.Grid = chart.GridLines
			}

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
				variables := constants
				variables["x"] = x
				y, _ := expression.Evaluate(variables)
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
