package math

import (
	"os"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
)

func LatexCommand() lib.Command {
	return lib.Command{
		Name:        "latex",
		Description: "Rendre du code LaTeX sur une image haute définition et l'envoyer dans le chat Telegram.",
		// TODO: customize output, documentation on installing imagick and pdflatex, documentation on available packages
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, _ map[string]interface{}) (err error) {
			expression := strings.Join(args, " ")

			path, file, err := lib.GenerateLatex(strconv.Itoa(update.UpdateID), expression)
			if err != nil {
				return lib.Error(bot, update, "Erreur dans le traitement de l'expression.")
			}
			reader := telegram.FileReader{
				Name:   "expression.png",
				Reader: file,
				Size:   -1,
			}
			photo := telegram.NewPhotoUpload(update.Message.Chat.ID, reader)
			_, err = bot.Send(photo)
			os.Remove(path)
			return
		},
	}
}
