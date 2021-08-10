package math

import (
	"fmt"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/math/src"
)

func CalcCommand() lib.Command {
	src.GatherFunctions()
	return lib.Command{
		Name:        "calc",
		Description: fmt.Sprintf("Calculer une valeur mathématique à l'aide d'une expression.\n%s\n\n%s", src.DataDocumentation, src.CalcDisclaimer),
		Flags: map[string]lib.Flag{
			"sf":  {"Nombre de chiffres après la virgule", 2, nil},
			"sci": {"Activer la notation scientifique (0 ou 1)", false, nil},
		},
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				return bot.Help(chatID, "calc")
			}

			expression := strings.Join(args, " ")
			if err := src.CheckExpression(expression); err != nil {
				return bot.Error(chatID, err.Error())
			}

			result, err := src.Evaluate(expression, 1.0)
			if err != nil {
				return bot.Error(chatID, err.Error())
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
