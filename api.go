package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
	body, err := MakeRequest(map[string]string{
		"q": query,
	})
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return
}

func MakeRequest(params map[string]string) (body []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(EncodeURL(params))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}

func EncodeURL(params map[string]string) string {
	url, _ := url.Parse(API_URL)
	query := url.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	url.RawQuery = query.Encode()
	return url.String()
}
