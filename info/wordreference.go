package info

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/html"

	"github.com/theovidal/parcolar/lib"
)

// WordReferenceUrl is the endpoint to the translation dictionary
const WordReferenceUrl = "https://www.wordreference.com"

func WordReferenceCommand() lib.Command {
	return lib.Command{
		Name:        "translation",
		Description: fmt.Sprintf("Obtenir la traduction d'un terme ou d'une expression (WordReference). Passez la langue source en premier argument, la langue cible en deuxième et le terme ou l'expression à la suite.\n\nLes langues disponibles sont : %s.", wordReferenceLanguagesDoc),
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, _ map[string]interface{}) (err error) {
			if len(args) < 3 {
				help := telegram.NewMessage(update.Message.Chat.ID, WordReferenceCommand().Help())
				help.ParseMode = "Markdown"
				_, err = bot.Send(help)
				return
			}
			from := args[0]
			to := args[1]
			search := strings.Join(args[2:], " ")

			if _, exists := wordReferenceLanguages[from]; !exists {
				return lib.Error(bot, update, "La langue source est invalide. Veuillez choisir parmi "+wordReferenceLanguagesDoc)
			}

			if _, exists := wordReferenceLanguages[to]; !exists {
				return lib.Error(bot, update, "La langue cible est invalide. Veuillez choisir parmi "+wordReferenceLanguagesDoc)
			}

			response, err := http.Get(fmt.Sprintf("%s/%s%s/%s", WordReferenceUrl, from, to, search))
			if err != nil {
				return
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
				return
			}

			audioSources := document.Find("audio source")

			var audios []string
			audioSources.Each(func(index int, audio *goquery.Selection) {
				url, _ := audio.Attr("src")
				urlParts := strings.Split(url, "/")

				label := "Source"
				if len(urlParts) > 3 {
					label = strings.Join(urlParts[2:4], "/")
				}

				audios = append(audios, fmt.Sprintf("[%s](https://wordreference.com%s)", label, url))
			})

			messages := []string{
				fmt.Sprintf("*―――――― 📚 %s → %s ――――――*\n\n🔊 %s", search, strings.ToUpper(to), strings.Join(audios, ", ")),
			}
			var messageIndex int

			selection := document.Find("table.WRD tbody tr:not(.wrtopsection,.langHeader)")
			if selection.Length() == 0 {
				return lib.Error(bot, update, "Aucune traduction trouvée pour ce terme ou cette expression.")
			}

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
					return
				}
			}
			return
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
		} else if child.Data != "br" && child.Data != "a" {
			text += child.Data
		}
	}
	return
}

var wordReferenceLanguages = map[string]string{
	"fr": "🇫🇷",
	"en": "🇬🇧",
	"es": "🇪🇸",
	"it": "🇮🇹",
	"de": "🇩🇪",
	"nl": "🇳🇱",
	"sv": "🇸🇪",
	"pt": "🇵🇹",
	"ru": "🇷🇺",
	"pl": "🇵🇱",
	"cz": "🇨🇿",
	"gr": "🇬🇷",
	"tr": "🇹🇷",
	"zh": "🇨🇳",
	"ja": "🇯🇵",
	"ko": "🇰🇴",
	"ar": "🇸🇦",
}

var wordReferenceLanguagesDoc = func() string {
	var languages []string
	for lang, flag := range wordReferenceLanguages {
		languages = append(languages, flag+" "+lang)
	}
	return strings.Join(languages, ", ")
}()
