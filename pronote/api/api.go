package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/theovidal/bacbot/lib"
)

// Response holds the data from Pronote
type Response struct {
	// A list of errors that could occur during the request
	Errors []interface{}
	// The message associated to an error
	Message string
	// User token (if requested)
	Token string
	// Request data (if requested)
	Data Data
}

// Data stores possible data from the Pronote API
type Data struct {
	// Homework to do for the next days
	Homeworks []Homework
	// Lessons during the next days
	Timetable []Lesson
	// Lesson contents written for the passed days
	Contents Contents
}

// MakeRequest executes a GraphQL query to the Pronote API
func MakeRequest(query string) (Response, error) {
	request, _ := http.NewRequest(
		"POST",
		os.Getenv("PRONOTE_API")+"/graphql",
		strings.NewReader(query),
	)

	request.Header.Add("Content-Type", "application/json")

	for i := 0; i < 3; i++ {
		var result Response
		token := lib.Cache.Get(context.Background(), "token").Val()
		request.Header.Set("Token", token)

		response, err := lib.DoRequest(request)
		if err != nil {
			return Response{}, err
		}
		var bytes []byte
		bytes, _ = ioutil.ReadAll(response.Body)
		response.Body.Close()

		_ = json.Unmarshal(bytes, &result)

		fmt.Println(response.Request.Host, response.StatusCode, string(bytes))

		if response.StatusCode == 200 && len(result.Errors) == 0 && result.Message == "" {
			return result, nil
		}

		err = Login()
		if err != nil {
			return Response{}, err
		}
	}

	return Response{}, errors.New("unable to make request after 3 trials")
}

// Login uses user's credentials to get an API token
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
