package info

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/html"

	"github.com/theovidal/bacbot/lib"
)

// WordReferenceUrl is the endpoint to the translation dictionary
const WordReferenceUrl = "https://www.wordreference.com"

func WordReferenceCommand() lib.Command {
	return lib.Command{
		Name:        "translation",
		Description: "Obtenir la traduction d'un terme ou d'une expression (WordReference)",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, _ map[string]interface{}) error {
			if len(args) < 3 {
				return lib.Error(bot, update, "Indiquez les deux langues ainsi que le terme à traduire.")
			}
			from := args[0]
			to := args[1]
			search := strings.Join(args[2:], " ")

			response, err := http.Get(fmt.Sprintf("%s/%s%s/%s", WordReferenceUrl, from, to, search))
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			if response.StatusCode == 404 {
				return lib.Error(bot, update, "La combinaison de langues est inconnue.")
			}
			if response.StatusCode != 200 {
				return lib.Error(bot, update, "Une erreur inconnue s'est produite lors de la recherche dans le dictionnaire.")
			}

			document, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			selection := document.Find("table.WRD tbody tr:not(.wrtopsection,.langHeader)")
			if selection.Length() == 0 {
				return lib.Error(bot, update, "Aucune traduction trouvée pour ce terme ou cette expression.")
			}

			messages := []string{
				fmt.Sprintf("*―――――― 📚 %s → %s ――――――*", search, strings.ToUpper(to)),
			}
			var messageIndex int

			selection.Each(func(_ int, element *goquery.Selection) {
				var content string
				if _, exists := element.Attr("id"); exists {
					element.Children().Each(func(_ int, child *goquery.Selection) {
						text, pronouns := ExtractTranslationText(child.Nodes[0])
						word := fmt.Sprintf("*%s* _%s_", text, pronouns)

						if class, _ := child.Attr("class"); class == "FrWrd" || class == "FrEx" {
							content += fmt.Sprintf("\n\n• %s ", word)
						} else if class == "ToWrd" || class == "ToEx" {
							content += "\n→ " + word
						} else {
							content += child.Text()
						}
					})
				} else {
					if element.Children().Length() == 3 {
						text, pronouns := ExtractTranslationText(element.Children().Get(2))
						content += fmt.Sprintf("\n     *%s* _%s_", text, pronouns)
					} else {
						content += fmt.Sprintf("\n_%s_", element.Text())
					}
				}

				if len(messages[messageIndex]+content) > 2000 {
					messageIndex += 1
					messages = append(messages, "")
				}
				messages[messageIndex] += content
			})

			for _, content := range messages {
				msg := telegram.NewMessage(update.Message.Chat.ID, content)
				msg.ParseMode = "Markdown"
				_, err = bot.Send(msg)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// ExtractTranslationText navigates deeply into HTML nodes to extract useful text for translations
func ExtractTranslationText(node *html.Node) (text, pronouns string) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Data == "strong" || child.Data == "span" {
			t, p := ExtractTranslationText(child)
			text += t
			pronouns += p
		} else if child.Data == "em" && child.FirstChild != nil {
			pronouns += child.FirstChild.Data
		} else if child.Data != "br" {
			text += child.Data
		}
	}
	return
}
