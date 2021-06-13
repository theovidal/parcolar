package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
)

// ParcoursupUrl is the endpoint of the OpenData API
const ParcoursupUrl = "https://data.enseignementsup-recherche.gouv.fr/api/records/1.0/search/?dataset=fr-esr-parcoursup"

// ParcoursupResponse stores the result of a call to the OpenData API
type ParcoursupResponse struct {
	// The list of record retrieved
	Records []ParcoursupRecord
}

// ParcoursupRecord stores a Parcoursup file
type ParcoursupRecord struct {
	// The unique identifier of the record
	ID string `json:"recordid"`
	// Date when the record was published
	Timestamp string `json:"record_timestamp"`
	// Specific fields related to the record
	Fields struct {
		// Full name of the file
		Name string `json:"g_ea_lib_vx"`
		// The concerned course
		Course string `json:"form_lib_voe_acc"`
		// On what domain the course is specialized
		Specialization string `json:"fil_lib_voe_acc"`
		// More details on the course (required certificate...)
		CourseDetail string `json:"detail_forma"`

		// French "département" where the course is
		Department string `json:"dep_lib"`
		// French region where the course is
		Region string `json:"region_etab_aff"`

		// Web URL to the Parcoursup file
		Link string `json:"lien_form_psup"`
	}
}

// SearchParcoursup queries the API to search for courses on Parcoursup
func SearchParcoursup(query string) (result ParcoursupResponse) {
	request, _ := http.NewRequest(
		"GET",
		lib.EncodeURL(ParcoursupUrl, map[string]string{
			"q": query,
		}),
		nil,
	)
	response, err := lib.DoRequest(request)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	response.Body.Close()

	err = json.Unmarshal(body, &result)
	return
}

// ParcoursupCommand processes an inline query from a user and returns the results for them to choose from
func ParcoursupCommand(bot *telegram.BotAPI, update *telegram.Update) {
	records := SearchParcoursup(update.InlineQuery.Query).Records

	var results []interface{}
	for _, record := range records {
		messageContent := fmt.Sprintf("[%s](%s)\n%s - %s", record.Fields.Name, record.Fields.Link, record.Fields.Course, record.Fields.Specialization)
		if record.Fields.CourseDetail != "" {
			messageContent += " - " + record.Fields.CourseDetail
		}
		messageContent += fmt.Sprintf("\n_%s, %s_", record.Fields.Department, record.Fields.Region)

		results = append(results, telegram.InlineQueryResultArticle{
			Type:  "article",
			ID:    record.ID,
			Title: fmt.Sprintf("%s - %s - %s", record.Fields.Name, record.Fields.Course, record.Fields.Specialization),
			URL:   record.Fields.Link,
			InputMessageContent: telegram.InputTextMessageContent{
				Text:                  messageContent,
				ParseMode:             "Markdown",
				DisableWebPagePreview: false,
			},
		})
	}

	_, err := bot.AnswerInlineQuery(telegram.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
	})
	if err != nil {
		log.Println(err)
	}
}
