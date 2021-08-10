package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/lib"
)

// ParcoursupUrl is the endpoint of the OpenData API
const ParcoursupUrl = "https://data.enseignementsup-recherche.gouv.fr/api/records/1.0/search/?dataset=fr-esr-cartographie_formations_parcoursup"

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
		Name string `json:"etab_nom"`
		// The concerned course
		Course string `json:"nm"`
		// On what domain the course is specialized
		Specialization string `json:"fl"`

		// French city where the course is
		City string `json:"commune"`
		// French "département" where the course is
		Department string `json:"departement"`
		// French region where the course is
		Region string `json:"region"`
		// Longitude and latitude of course's location
		Coordinates []float64 `json:"etab_gps"`

		// Web URL to the Parcoursup file
		File string `json:"fiche"`
		// Web URL to course's website
		Website string `json:"etab_url"`
		// Web URL to last year statistics on this course
		Statistics string `json:"dataviz"`
	}
}

// SearchParcoursup queries the API to search for courses on Parcoursup
func SearchParcoursup(query string) (result ParcoursupResponse, err error) {
	request, _ := http.NewRequest(
		"GET",
		lib.EncodeURL(ParcoursupUrl, map[string]string{
			"q": query,
		}),
		nil,
	)
	response, err := http.DefaultClient.Do(request)
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
func ParcoursupCommand(bot *lib.Bot, update *telegram.Update) (err error) {
	response, err := SearchParcoursup(update.InlineQuery.Query)
	if err != nil {
		return
	}
	records := response.Records

	var results []interface{}
	for _, record := range records {
		markup := &telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				{
					{Text: fmt.Sprintf("%s — %s", record.Fields.Course, record.Fields.Name), URL: &record.Fields.Website},
				},
				{
					{Text: "📜 Fiche détaillée sur Parcoursup", URL: &record.Fields.File},
				},
			},
		}
		/* TODO: URL not added to button for some requests
		if record.Fields.Statistics != "" {
			markup.InlineKeyboard = append(markup.InlineKeyboard, []telegram.InlineKeyboardButton{
				{ Text: "📈 Données statistiques pour l'année antérieure", URL: &record.Fields.Statistics },
			})
		} */

		results = append(results, telegram.InlineQueryResultLocation{
			Type:        "location",
			ID:          record.ID,
			Title:       fmt.Sprintf("%s - %s - %s", record.Fields.Course, record.Fields.Name, record.Fields.Specialization),
			Latitude:    record.Fields.Coordinates[0],
			Longitude:   record.Fields.Coordinates[1],
			ReplyMarkup: markup,
		})
	}

	_, err = bot.AnswerInlineQuery(telegram.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       results,
	})
	return
}
