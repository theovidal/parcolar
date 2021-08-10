package wolfram

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
)

func Command() lib.Command {
	return lib.Command{
		Name:        "wolfram",
		Description: "Réaliser une requête en utilisant la version en ligne de Wolfram|Alpha et compte gratuit : par conséquent, certaines fonctionnalités ne seront disponibles. L'outil est efficace pour toutes recherches générales, et les calculs mathématiques sans étapes.",
		Execute: func(bot *lib.Bot, update *telegram.Update, chatID int64, args []string, _ map[string]interface{}) (err error) {
			if len(args) == 0 {
				return bot.Help(chatID, "wolfram")
			}

			query := strings.Join(args, " ")
			result, err := makeRequest(query)
			if err != nil {
				lib.LogError(err.Error())
				return bot.Error(chatID, "Une erreur est survenue lors de la requête.")
			}

			if result.HasSuccess && !result.HasError {
				var photoGroups [][]interface{}
				groupIndex := 0
				photoGroups = append(photoGroups, []interface{}{})

				for _, pod := range result.Pods {
					for _, subpod := range pod.Subpods {
						photo := telegram.NewInputMediaPhoto(subpod.Image.URL)
						photo.Caption = pod.Title
						if len(photoGroups[groupIndex]) == 10 {
							groupIndex++
							photoGroups = append(photoGroups, []interface{}{})
						}
						photoGroups[groupIndex] = append(photoGroups[groupIndex], photo)
					}
				}

				for _, group := range photoGroups {
					var media telegram.Chattable = telegram.NewMediaGroup(chatID, group)
					if len(group) == 1 {
						media = telegram.NewPhotoUpload(chatID, group[0])
					}

					if _, err = bot.Send(media); err != nil {
						lib.LogError(err.Error())
						return bot.Error(chatID, "Une erreur est survenue lors de l'envoi des photos.")
					}
				}
			}

			if !result.HasSuccess && !result.HasError {
				message := telegram.NewMessage(chatID, "")
				for _, tip := range result.Tips.Data {
					message.Text += fmt.Sprintf("💡 %s\n", tip)
				}

				if message.Text == "" {
					message.Text = "La requête n'a retourné aucun résultat."
				}
				message.Text += "\n"

				var buttons [][]telegram.InlineKeyboardButton
				for _, didyoumean := range result.DidYouMeans.Data {
					callback := "/wolfram " + didyoumean
					text := fmt.Sprintf("❓ Souhaitez-vous dire: %s ?", didyoumean)
					if len(callback) > 64 {
						message.Text += text + "\n"
					} else {
						buttons = append(buttons, []telegram.InlineKeyboardButton{
							{
								Text:         text,
								CallbackData: &callback,
							},
						})
					}
				}
				message.ReplyMarkup = telegram.InlineKeyboardMarkup{
					InlineKeyboard: buttons,
				}
				_, err = bot.Send(message)
			}

			if !result.HasSuccess && result.HasError {
				return bot.Error(chatID, fmt.Sprintf("Erreurs sur la requête : `%s`.", result.Error.Content))
			}

			return
		},
	}
}
