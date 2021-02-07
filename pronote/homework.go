package pronote

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HomeworkCommand(bot *telegram.BotAPI, update telegram.Update, _ []string) error {
	response, err := GetHomework()
	if err != nil {
		return err
	}

	if len(response.Homeworks) == 0 {
		msg := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun devoir n'a été rédigé pour le moment.")
		_, err = bot.Send(msg)
		return err
	}

	content := ""
	for _, homework := range response.Homeworks {
		content += homework.String()
	}

	msg := telegram.NewMessage(update.Message.Chat.ID, content)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	_, err = bot.Send(msg)

	return err
}
