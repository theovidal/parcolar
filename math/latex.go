package math

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/math/data"
)

func LatexCommand() lib.Command {
	return lib.Command{
		Name:        "latex",
		Description: "Rendre du code LaTeX sur une image haute définition et l'envoyer dans le chat Telegram.",
		Flags: map[string]lib.Flag{
			"background": {"Couleur de l'arrière-plan", "white", &data.Colors},
			"text":       {"Couleur du texte", "black", &data.Colors},
		},
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, chatID int64, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				return lib.Help(bot, chatID, LatexCommand())
			}

			backgroundColor, textColor := flags["background"].(string), flags["text"].(string)
			expression := fmt.Sprintf("\\color{%s} $%s$", textColor, strings.Join(args, " "))

			path, photo, err := lib.GenerateLatex(
				strconv.Itoa(update.UpdateID),
				fmt.Sprintf("\\pagecolor{%s}", backgroundColor),
				expression,
			)
			if err != nil {
				return lib.Error(bot, chatID, err.Error())
			}
			photoReader := telegram.FileReader{
				Name:   "expression.png",
				Reader: photo,
				Size:   -1,
			}
			photoUpload := telegram.NewPhotoUpload(chatID, photoReader)
			_, err = bot.Send(photoUpload)
			os.Remove(path)
			return
		},
	}
}
