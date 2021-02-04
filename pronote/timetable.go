package pronote

import (
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func TimetableCommand(bot *telegram.BotAPI, update telegram.Update, _ []string) error {
	response, err := GetTimetable(false)
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

func TimetableTicker(bot *telegram.BotAPI) error {
	response, err := GetTimetable(false)
	if err != nil {
		return err
	}

	if len(response.Data.Timetable) == 0 {
		return nil
	}

	nextLesson := response.Data.Timetable[0]

	for _, lesson := range response.Data.Timetable {
		date := int64(lesson.From) / 1000
		from := time.Now().Unix()
		to := time.Now().Add(time.Minute * 10).Unix()

		if nextLesson.Cancelled || nextLesson.Status == "Prof. absent" || date > to {
			break
		}

		if date >= from && date <= to {
			content := "🔔 *Prochain cours*\n" + lesson.String()
			msg := telegram.NewMessage(663102119, content)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			return err
		}
	}

	return nil
}
