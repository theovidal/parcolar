package pronote

import (
	"fmt"
	"os"
	"strconv"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"

	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote/api"
)

// TimetableLoop runs the TimetableTicker every 10 minutes, and is meant to be used in a goroutine.
func TimetableLoop(bot *telegram.BotAPI) {
	for range time.Tick(time.Minute * 10) {
		if err := TimetableTicker(bot); err != nil {
			lib.LogError("‼ Error handling timetable ticker: %s", err)
		}
	}
}

// TimetableTicker periodically fetches the timetable on PRONOTE for upcoming lessons, and sends a notification if there is one in the next 10 minutes
func TimetableTicker(bot *telegram.BotAPI) (err error) {
	response, err := api.GetTimetable(now.BeginningOfDay(), now.BeginningOfDay().Add(time.Hour*26))
	if err != nil || len(response.Timetable) == 0 {
		return
	}

	location, err := time.LoadLocation(os.Getenv("PRONOTE_TIMEZONE"))
	if err != nil {
		return fmt.Errorf("Can't get timezone defined in PRONOTE_TIMEZONE environment variable")
	}

	from := time.Now().In(location).Unix()
	to := time.Now().In(location).Add(time.Minute * 10).Unix()

	for _, lesson := range response.Timetable {
		date := int64(lesson.From) / 1000

		if date > to {
			break
		}

		if lesson.Cancelled || lesson.Status == "Prof. absent" || lesson.Remote {
			continue
		}

		if date >= from && date <= to {
			content := "*―――――― 🔔 Prochain cours ――――――*\n" + lesson.String()
			chatID, _ := strconv.Atoi(os.Getenv("TELEGRAM_CHAT"))
			message := telegram.NewMessage(int64(chatID), content)
			message.ParseMode = "MarkdownV2"
			_, err = bot.Send(message)
			return
		}
	}

	return
}
