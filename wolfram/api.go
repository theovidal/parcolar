package wolfram

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/theovidal/parcolar/lib"
)

const ApiUrl = "https://api.wolframalpha.com/v2/query"

func makeRequest(query string) (result QueryResult, err error) {
	url := lib.EncodeURL(ApiUrl, map[string]string{
		"appid": os.Getenv("WOLFRAM_ID"),
		"input": query,
	})
	response, err := http.Get(url)
	if err != nil {
		return
	}
	var body []byte
	body, _ = ioutil.ReadAll(response.Body)
	response.Body.Close()

	err = xml.Unmarshal(body, &result)
	return
}
