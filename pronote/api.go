package pronote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/theovidal/bacbot/lib"
)

type Response struct {
	Errors  []interface{}
	Message string
	Token   string
	Data    PronoteData
}

type PronoteData struct {
	Homeworks []Homework
	Timetable []Lesson
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

func GetHomework() (PronoteData, error) {
	query := ParseGraphql(fmt.Sprintf(`
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

	response, err := MakeRequest(query)
	return response.Data, err
}

func GetTimetable(todayOnly bool) (PronoteData, error) {
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

func MakeRequest(query string) (result Response, err error) {
	request, _ := http.NewRequest(
		"POST",
		os.Getenv("PRONOTE_API")+"/graphql",
		strings.NewReader(query),
	)

	request.Header.Add("Content-Type", "application/json")

	var response *http.Response
	for {
		var currentResult Response
		token := lib.Cache.Get(context.Background(), "token").Val()
		request.Header.Set("Token", token)

		response, err = lib.DoRequest(request)
		if err != nil {
			return
		}
		var bytes []byte
		bytes, _ = ioutil.ReadAll(response.Body)
		response.Body.Close()

		_ = json.Unmarshal(bytes, &currentResult)

		fmt.Println(response.Request.Host, response.StatusCode, string(bytes))

		fmt.Println(len(currentResult.Errors), currentResult.Message)
		if response.StatusCode == 200 && len(currentResult.Errors) == 0 && currentResult.Message == "" {
			result = currentResult
			break
		}

		err = Login()
		if err != nil {
			return
		}
	}

	return
}

func Login() error {
	query, _ := json.Marshal(map[string]string{
		"url":      os.Getenv("PRONOTE_SERVER"),
		"cas":      os.Getenv("PRONOTE_CAS"),
		"username": os.Getenv("PRONOTE_USER"),
		"password": os.Getenv("PRONOTE_PASSWORD"),
	})

	request, _ := http.NewRequest(
		"POST",
		os.Getenv("PRONOTE_API")+"/auth/login",
		bytes.NewReader(query),
	)

	request.Header.Add("Content-Type", "application/json")
	response, err := lib.DoRequest(request)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, _ = ioutil.ReadAll(response.Body)
	response.Body.Close()

	var result Response
	_ = json.Unmarshal(bytes, &result)

	lib.Cache.Set(context.Background(), "token", result.Token, 0)

	return nil
}

func ParseGraphql(query string) string {
	raw, _ := json.Marshal(map[string]string{
		"query": query,
	})
	return string(raw)
}
