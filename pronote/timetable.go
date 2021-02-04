package pronote

import (
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func TimetableCommand(bot *telegram.BotAPI, update telegram.Update, _ []string) error {
	response, err := GetTimetable()
	if err != nil {
		return err
	}

	days := make(map[string]string)
	for _, lesson := range response.Data.Timetable {
		day := time.Unix(int64(lesson.From/1000), 0).Format("02/01")

		days[day] = days[day] + lesson.String()
	}

	var content string
	for day, lessons := range days {
		content += fmt.Sprintf("*―――――― %s ――――――*\n%s\n", day, lessons)
	}

	msg := telegram.NewMessage(update.Message.Chat.ID, content)
	msg.ParseMode = "MarkdownV2"
	_, err = bot.Send(msg)

	return err
}
