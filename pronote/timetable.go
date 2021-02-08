package pronote

import (
	"fmt"
	"os"
	"strconv"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/pronote/api"
)

func TimetableCommand(bot *telegram.BotAPI, update telegram.Update, _ []string) error {
	response, err := api.GetTimetable(false)
	if err != nil {
		return err
	}

	if len(response.Timetable) == 0 {
		msg := telegram.NewMessage(update.Message.Chat.ID, "🍃 Aucun cours n'est prévu pour le moment.")
		_, err = bot.Send(msg)
		return err
	}

	days := make(map[string]string)
	for _, lesson := range response.Timetable {
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
	response, err := api.GetTimetable(false)
	if err != nil {
		return err
	}

	if len(response.Timetable) == 0 {
		return nil
	}

	nextLesson := response.Timetable[0]

	for _, lesson := range response.Timetable {
		date := int64(lesson.From) / 1000
		from := time.Now().Unix()
		to := time.Now().Add(time.Minute * 10).Unix()
		// from := int64(1613976780)
		// to := int64(1613977000)
		fmt.Println(nextLesson.From, from, to)

		if nextLesson.Cancelled || nextLesson.Status == "Prof. absent" || date > to {
			break
		}

		if date >= from && date <= to {
			content := "*―――――― 🔔 Prochain cours ――――――*\n" + lesson.String()
			chat, _ := strconv.Atoi(os.Getenv("TELEGRAM_CHAT"))
			msg := telegram.NewMessage(int64(chat), content)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			return err
		}
	}

	return nil
}
