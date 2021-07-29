package math

import (
	"fmt"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/math/data"
)

func CalcCommand() lib.Command {
	data.GatherFunctions()
	return lib.Command{
		Name:        "calc",
		Description: fmt.Sprintf("Calculer une valeur mathématique à l'aide d'une expression.\n%s\n\n%s", data.DataDocumentation, data.CalcDisclaimer),
		Flags: map[string]lib.Flag{
			"sf":  {"Nombre de chiffres après la virgule", 2, nil},
			"sci": {"Activer la notation scientifique (0 ou 1)", false, nil},
		},
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				return lib.Help(bot, chatID, CalcCommand())
			}

			expression := strings.Join(args, " ")
			if err := data.CheckExpression(expression); err != nil {
				return lib.Error(bot, chatID, err.Error())
			}

			result, err := data.Evaluate(expression, 1.0)
			if err != nil {
				return lib.Error(bot, chatID, err.Error())
			}

			format := "= %." + strconv.Itoa(flags["sf"].(int))
			if flags["sci"].(bool) {
				format += "e"
			} else {
				format += "f"
			}
			message := telegram.NewMessage(chatID, fmt.Sprintf(format, result))
			_, err = bot.Send(message)
			return
		},
	}
}
