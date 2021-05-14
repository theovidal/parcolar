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

			res, err := http.Get(fmt.Sprintf("https://www.wordreference.com/%s%s/%s", from, to, search))
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode == 404 {
				return lib.Error(bot, update, "La combinaison de langues est inconnue.")
			}
			if res.StatusCode != 200 {
				return lib.Error(bot, update, "Une erreur inconnue s'est produite lors de la recherche dans le dictionnaire.")
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			selection := doc.Find("table.WRD tbody tr:not(.wrtopsection,.langHeader)")
			if selection.Length() == 0 {
				return lib.Error(bot, update, "Aucune traduction trouvée pour ce terme ou cette expression.")
			}

			parts := []string{
				fmt.Sprintf("*―――――― 📚 %s → %s ――――――*", search, strings.ToUpper(to)),
			}
			var part int

			selection.Each(func(_ int, element *goquery.Selection) {
				var content string
				if _, exists := element.Attr("id"); exists {
					element.Children().Each(func(_ int, child *goquery.Selection) {
						text, pronouns := ExtractText(child.Nodes[0])
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
						text, pronouns := ExtractText(element.Children().Get(2))
						content += fmt.Sprintf("\n     *%s* _%s_", text, pronouns)
					} else {
						content += fmt.Sprintf("\n_%s_", element.Text())
					}
				}

				if len(parts[part]+content) > 2000 {
					part += 1
					parts = append(parts, "")
				}
				parts[part] += content
			})

			for _, content := range parts {
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

func ExtractText(node *html.Node) (text, pronouns string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "strong" || c.Data == "span" {
			t, p := ExtractText(c)
			text += t
			pronouns += p
		} else if c.Data == "em" && c.FirstChild != nil {
			pronouns += c.FirstChild.Data
		} else if c.Data != "br" {
			text += c.Data
		}
	}
	return
}
