package wolfram

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/theovidal/parcolar/lib"
)

func Command() lib.Command {
	return lib.Command{
		Name: "wolfram",
		Description: "Réaliser une recherche sur Wolfram|Alpha (compte gratuit)",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, _ map[string]interface{}) (err error) {
			if len(args) == 0 {
				help := telegram.NewMessage(update.Message.Chat.ID, Command().Help())
				help.ParseMode = "Markdown"
				_, err = bot.Send(help)
				return
			}

			query := strings.Join(args, " ")
			result, err := makeRequest(query)
			if err != nil {
				lib.LogError(err.Error())
				return lib.Error(bot, update, "Une erreur est survenue lors de la requête.")
			}

			if result.HasSuccess && !result.HasError {
				for _, pod := range result.Pods {
					var photos []interface{}
					for _, subpod := range pod.Subpods {
						photo := telegram.NewInputMediaPhoto(subpod.Image.URL)
						photo.Caption = fmt.Sprintf("%s — %s", pod.Title, subpod.Text)
						photos = append(photos, photo)
					}
					if _, err = bot.Send(telegram.NewMediaGroup(update.Message.Chat.ID, photos)); err != nil {
						lib.LogError(err.Error())
						return lib.Error(bot, update, "Une erreur est survenue lors de l'envoi des photos.")
					}
				}
			}

			if !result.HasSuccess && !result.HasError {
				var message string
				for _, tip := range result.Tips.Data {
					message += fmt.Sprintf("💡 %s\n", tip)
				}
				message += "\n"
				for _, didyoumean := range result.DidYouMeans.Data {
					message += fmt.Sprintf("❓ Souhaitez-vous dire: %s ?\n", didyoumean)
				}
				_, err = bot.Send(telegram.NewMessage(update.Message.Chat.ID, message))
			}

			if !result.HasSuccess && result.HasError {
				return lib.Error(bot, update, fmt.Sprintf("Erreurs sur la requête : `%s`.", result.Error.Content))
			}

			return
		},
	}
}
