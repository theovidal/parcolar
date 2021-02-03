package lib

import (
	"net/http"
	"net/url"
)

func MakeRequest(request *http.Request) (response *http.Response, err error) {
	response, err = http.DefaultClient.Do(request)
	return
}

func EncodeURL(input string, params map[string]string) string {
	url, _ := url.Parse(input)
	query := url.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	url.RawQuery = query.Encode()
	return url.String()
}
