package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/pronote/api"
)

func ContentsCommand(bot *telegram.BotAPI, update telegram.Update, _ []string) error {
	response, err := api.GetContents()
	if err != nil {
		return err
	}

	if len(response.Contents) == 0 {
		msg := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun contenu de cours n'a été saisi pour le moment.")
		_, err = bot.Send(msg)
		return err
	}

	output := ""
	for _, content := range response.Contents.Reverse() {
		output += content.String()
	}

	msg := telegram.NewMessage(update.Message.Chat.ID, output)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	_, err = bot.Send(msg)

	return err
}
