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
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) (err error) {
			if len(args) == 0 {
				help := telegram.NewMessage(update.Message.Chat.ID, LatexCommand().Help())
				help.ParseMode = "Markdown"
				_, err := bot.Send(help)
				return err
			}

			bgColor := flags["background"].(string)
			textColor := flags["text"].(string)

			expression := fmt.Sprintf("\\color{%s} $%s$", textColor, strings.Join(args, " "))

			path, file, err := lib.GenerateLatex(
				strconv.Itoa(update.UpdateID),
				fmt.Sprintf("\\pagecolor{%s}", bgColor),
				expression,
			)
			if err != nil {
				return lib.Error(bot, update, err.Error())
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
