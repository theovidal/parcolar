package parcoursup

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/theovidal/bacbot/lib"
)

// ApiUrl is the endpoint of the OpenData API
const ApiUrl = "https://data.enseignementsup-recherche.gouv.fr/api/records/1.0/search/?dataset=fr-esr-parcoursup"

// APIResult stores the result of a call to the OpenData API
type APIResult struct {
	// The list of record retrieved
	Records []Record
}

// Record stores a Parcoursup file
type Record struct {
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

// SearchRecords queries the API to search for courses on Parcoursup
func SearchRecords(query string) (result APIResult) {
	request, _ := http.NewRequest(
		"GET",
		lib.EncodeURL(ApiUrl, map[string]string{
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
