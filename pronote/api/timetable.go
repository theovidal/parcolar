package api

import (
	"fmt"
	"time"

	"github.com/theovidal/bacbot/lib"
)

type Lesson struct {
	From      int
	Subject   string
	Teacher   string
	Room      string
	Status    string
	Cancelled bool `json:"isCancelled"`
	Remote    bool `json:"remoteLesson"`
}

func (lesson *Lesson) String() (output string) {
	var emoji string
	if lesson.Remote {
		emoji = "🏡"
	} else if lesson.Cancelled || lesson.Status == "Prof. absent" {
		emoji = "🚪"
	} else {
		emoji = "🕑"
	}
	output += fmt.Sprintf(
		"%s *%s: %s* \\(%s\\)",
		emoji,
		time.Unix(int64(lesson.From/1000), 0).Format("15h04"),
		lib.ParseTelegramMessage(Subjects[lesson.Subject].Name),
		lib.ParseTelegramMessage(lesson.Teacher),
	)
	if lesson.Room != "" {
		output += " — " + lib.ParseTelegramMessage(lesson.Room)
	}

	if lesson.Cancelled || lesson.Status == "Prof. absent" {
		output = "~" + output
		output += "~"
	}
	if lesson.Status != "" {
		output += fmt.Sprintf(" `%s`", lesson.Status)
	}
	output += "\n"
	return
}

func GetTimetable(todayOnly bool) (Data, error) {
	from := time.Now().Format("2006-01-02")
	toTime := time.Now().Add(time.Hour * 24)
	if todayOnly {
		toTime.Add(time.Hour * 24 * 6)
	}
	to := toTime.Format("2006-01-02")

	query := ParseGraphql(fmt.Sprintf(`
		{
			timetable(from: "%s", to: "%s") {
				from
				to
				subject
				teacher
				room
				status
				isCancelled
				remoteLesson
			}
		}
	`, from, to),
	)

	response, err := MakeRequest(query)
	return response.Data, err
}
