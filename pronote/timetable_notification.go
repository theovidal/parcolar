package pronote

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"

	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/pronote/api"
)

// TimetableLoop runs the TimetableTicker every 10 minutes, and is meant to be used in a goroutine.
func TimetableLoop(bot *telegram.BotAPI) {
	for range time.Tick(time.Minute * 10) {
		err := TimetableTicker(bot)
		if err != nil {
			log.Println(lib.Red.Sprintf("‼ Error handling timetable ticker: %s", err))
		}
	}
}

// TimetableTicker periodically fetches the timetable on PRONOTE for upcoming lessons, and sends a notification if there is one in the next 10 minutes
func TimetableTicker(bot *telegram.BotAPI) error {
	response, err := api.GetTimetable(now.BeginningOfDay(), now.BeginningOfDay().Add(time.Hour*26))
	if err != nil {
		return err
	}

	if len(response.Timetable) == 0 {
		return nil
	}

	from := time.Now().Unix()
	to := time.Now().Add(time.Minute * 10).Unix()

	for _, lesson := range response.Timetable {
		location, err := time.LoadLocation(os.Getenv("PRONOTE_TIMEZONE"))
		if err != nil {
			return errors.New("Can't get timezone " + err.Error() + " (from PRONOTE_TIMEZONE environment variable)")
		}
		date := time.Unix(int64(lesson.From)/1000, 0).In(location).Unix()

		if date > to {
			break
		}

		if lesson.Cancelled || lesson.Status == "Prof. absent" || lesson.Remote {
			continue
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
