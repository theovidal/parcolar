package lib

import (
	"net/http"
	"net/url"
)

func DoRequest(request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
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
