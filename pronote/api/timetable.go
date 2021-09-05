package api

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/theovidal/parcolar/lib"
)

// Lesson represents a class during a day on a specific subject
type Lesson struct {
	From      int
	To        int
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
		"%s *%s\\-%s: %s* \\(%s\\)",
		emoji,
		time.Unix(int64(lesson.From/1000), 0).Format("15h04"),
		time.Unix(int64(lesson.To/1000), 0).Format("15h04"),
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

// GetTimetable fetches the timetable for the next 7 days
func GetTimetable(cache *redis.Client, from time.Time, to time.Time) (Data, error) {
	query := ParseGraphQL(fmt.Sprintf(`
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
	`, from.Format("2006-01-02"), to.Format("2006-01-02")),
	)

	response, err := MakeRequest(cache, query)
	return response.Data, err
}
