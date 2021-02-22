package math

import (
	"fmt"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
)

func CalcCommand() lib.Command {
	return lib.Command{
		Name:        "calc",
		Description: fmt.Sprintf("Calculer une valeur mathématique à l'aide d'une expression.\n%s\n\n%s", dataDocumentation, calcDisclaimer),
		Flags: map[string]lib.Flag{
			"sf":  {"Nombre de chiffres après la virgule", 2},
			"sci": {"Activer la notation scientifique (0 ou 1)", 0},
		},
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error {
			if len(args) == 0 {
				help := telegram.NewMessage(update.Message.Chat.ID, CalcCommand().Help())
				help.ParseMode = "Markdown"
				_, err := bot.Send(help)
				return err
			}

			function := strings.Join(args, " ")
			msg := CheckExpression(function)
			if msg != "" {
				return lib.Error(bot, update, msg)
			}

			value := Evaluate(function, 1.0)

			format := "= %." + strconv.Itoa(flags["sf"].(int))
			if flags["sci"].(int) == 1 {
				format += "e"
			} else {
				format += "f"
			}
			_, err := bot.Send(telegram.NewMessage(update.Message.Chat.ID, fmt.Sprintf(format, value)))
			return err
		},
	}
}
