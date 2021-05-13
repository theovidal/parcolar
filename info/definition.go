package info

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
)

func DefinitionCommand() lib.Command {
	return lib.Command{
		Name:        "definition",
		Description: "Obtenir la définition d'un terme dans le dictionnaire (Larousse)",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error {
			if len(args) == 0 {
				return lib.Error(bot, update, "Merci d'indiquer un terme pour en chercher la définition dans le dictionnaire.")
			}

			word := args[0]
			res, err := http.Get("https://www.larousse.fr/dictionnaires/francais/" + word)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				return lib.Error(bot, update, "Une erreur inconnue s'est produite lors de la recherche dans le dictionnaire.")
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			selection := doc.Find("ul.Definitions li")
			if selection.Length() == 0 {
				return lib.Error(bot, update, "Aucune définition trouvée pour ce terme.")
			}

			content := fmt.Sprintf("*―――――― 📜 %s ――――――*\n", strings.ToUpper(word))
			selection.Each(func(_ int, definition *goquery.Selection) {
				content += fmt.Sprintf("\n• %s", definition.Text())
			})

			msg := telegram.NewMessage(update.Message.Chat.ID, content)
			msg.ParseMode = "Markdown"
			_, err = bot.Send(msg)
			return err
		},
	}
}
