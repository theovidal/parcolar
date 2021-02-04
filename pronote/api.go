package pronote

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/theovidal/parcoursupbot/lib"
)

type Response struct {
	Data struct {
		Homeworks []Homework
		Timetable []Lesson
	}
}

type Homework struct {
	Description string
	Subject     string
	Due         int `json:"for"`
	Files       []struct {
		Name string
		URL  string
	}
}

func (homework *Homework) String() (output string) {
	subject := subjects[homework.Subject]
	output += fmt.Sprintf(
		"*%s %s — %s*\n%s",
		subject.Emoji,
		subject.Name,
		time.Unix(int64(homework.Due/1000), 0).Format("02/01"),
		homework.Description,
	)
	for _, file := range homework.Files {
		output += fmt.Sprintf("\n📎 [%s](%s)", file.Name, file.URL)
	}
	output += "\n\n"
	return
}

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
		"%s *%s: %s* \\(%s\\)",
		emoji,
		time.Unix(int64(lesson.From/1000), 0).Format("15h04"),
		lib.ParseTelegramMessage(subjects[lesson.Subject].Name),
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

func GetHomework() (result Response, err error) {
	body := ParseGraphql(fmt.Sprintf(`
		{
			homeworks(from: "%s", to: "%s") {
				description
				subject
				for
				files {
					name
					url
				}
			}
		}
	`, time.Now().Format("2006-01-02"), time.Now().Add(time.Hour*24*15).Format("2006-01-02")),
	)

	content, err := MakeRequest("graphql", body)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(content), &result)
	return
}

func GetTimetable(todayOnly bool) (result Response, err error) {
	from := time.Now().Format("2006-01-02")
	var to string
	if todayOnly {
		to = from
	} else {
		to = time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02")
	}

	body := ParseGraphql(fmt.Sprintf(`
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

	content, err := MakeRequest("graphql", body)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(content), &result)
	return
}

func MakeRequest(endpoint string, query string) (content string, err error) {
	request, _ := http.NewRequest(
		"POST",
		os.Getenv("PRONOTE_API")+"/"+endpoint,
		strings.NewReader(query),
	)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Token", lib.Cache.Get(context.Background(), "token").Val())

	var response *http.Response
	for {
		response, err = lib.MakeRequest(request)
		if err != nil {
			return
		}

		if response.StatusCode != 500 {
			break
		}

		err = Login()
		if err != nil {
			return
		}
	}

	var bytes []byte
	bytes, err = ioutil.ReadAll(response.Body)
	response.Body.Close()

	content = string(bytes)
	return
}

func Login() error {
	query, _ := json.Marshal(map[string]string{
		"url":      os.Getenv("PRONOTE_SERVER"),
		"cas":      os.Getenv("PRONOTE_CAS"),
		"username": os.Getenv("PRONOTE_USER"),
		"password": os.Getenv("PRONOTE_PASSWORD"),
	})
	content, err := MakeRequest("auth/login", string(query))
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	json.Unmarshal([]byte(content), &data)

	lib.Cache.Set(context.Background(), "token", data["token"].(string), 0)

	return nil
}

func ParseGraphql(query string) string {
	raw, _ := json.Marshal(map[string]string{
		"query": query,
	})
	return string(raw)
}
