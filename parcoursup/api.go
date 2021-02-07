package parcoursup

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/theovidal/bacbot/lib"
)

var API_URL = "https://data.enseignementsup-recherche.gouv.fr/api/records/1.0/search/?dataset=fr-esr-parcoursup"

type APIResult struct {
	Records []Record
}

type Record struct {
	ID        string `json:"recordid"`
	Timestamp string `json:"record_timestamp"`
	Fields    struct {
		Name           string `json:"g_ea_lib_vx"`
		Course         string `json:"form_lib_voe_acc"`
		Specialization string `json:"fil_lib_voe_acc"`
		CourseDetail   string `json:"detail_forma"`

		Department string `json:"dep_lib"`
		Region     string `json:"region_etab_aff"`

		Link string `json:"lien_form_psup"`
	}
}

func SearchRecords(query string) (result APIResult) {
	request, _ := http.NewRequest(
		"GET",
		lib.EncodeURL(API_URL, map[string]string{
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
