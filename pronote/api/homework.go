package api

import (
	"fmt"
	"github.com/theovidal/bacbot/lib"
	"time"
)

// Homework defines the homework given for the next days by a teacher
type Homework struct {
	Description string
	Subject     string
	Due         int `json:"for"`
	Done        bool
	Files       []File
}

func (homework *Homework) String() (output string) {
	subject := Subjects[homework.Subject]
	var emoji string
	if homework.Done {
		emoji = "✅"
	} else {
		emoji = subject.Emoji
	}

	output += fmt.Sprintf(
		"*%s %s — %s*\n%s",
		emoji,
		lib.ParseTelegramMessage(subject.Name),
		time.Unix(int64(homework.Due/1000), 0).Format("02/01"),
		lib.ParseTelegramMessage(homework.Description),
	)
	for _, file := range homework.Files {
		output += file.String()
	}
	if homework.Done {
		output = "~" + output
		output += "~"
	}
	output += "\n\n"
	return
}

// GetHomework fetches the homework to do for the next 15 days
func GetHomework() (Data, error) {
	query := ParseGraphQL(fmt.Sprintf(`
		{
			homeworks(from: "%s", to: "%s") {
				description
				subject
				for
				done
				files {
					name
					url
				}
			}
		}
	`, time.Now().Format("2006-01-02"), time.Now().Add(time.Hour*24*15).Format("2006-01-02")),
	)

	response, err := MakeRequest(query)
	return response.Data, err
}
