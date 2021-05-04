package api

import (
	"fmt"
	"time"
)

// Content stores the content of a passed lesson
type Content struct {
	Subject  string
	Teachers string
	Time     int `json:"from"`
	Contents []ContentPart
}

// ContentPart stores the different parts of a lesson
type ContentPart struct {
	Name        string
	Category    string
	Description string
	Files       []File
}

func (content *Content) String() (output string) {
	subject := Subjects[content.Subject]
	output += fmt.Sprintf(
		"*%s %s — %s*\n",
		subject.Emoji,
		subject.Name,
		time.Unix(int64(content.Time/1000), 0).Format("02/01"),
	)

	for _, part := range content.Contents {
		if part.Category != "" {
			output += fmt.Sprintf("*« %s »* ", part.Category)
		}
		output += part.Description + "\n"
		for _, file := range part.Files {
			output += file.String()
		}
		output += "\n"
	}
	if len(content.Contents) == 0 {
		output += "\n"
	}
	return
}

// GetContents fetches lesson contents for the past 5 days
func GetContents() (Data, error) {
	query := ParseGraphQL(fmt.Sprintf(`
		{
			contents(from: "%s", to: "%s") {
				subject
				teachers
				from
				to
				color
				contents {
					name
					category
					htmlDescription
					description
					files {
						name
						url
					}
				}
			}
		}
	`, time.Now().AddDate(0, 0, -4).Format("2006-01-02"), time.Now().Format("2006-01-02")),
	)

	response, err := MakeRequest(query)
	return response.Data, err
}

// Contents is a shortcut for []Content
type Contents []Content

// Reverse sorts Content objects to reverse their order in a list
func (c Contents) Reverse() Contents {
	for i := 0; i < len(c)/2; i++ {
		j := len(c) - i - 1
		c[i], c[j] = c[j], c[i]
	}
	return c
}
