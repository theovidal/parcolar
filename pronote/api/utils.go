package api

import (
	"encoding/json"
	"fmt"
	"github.com/theovidal/bacbot/lib"
)

// File stores a document attached to homework or contents
type File struct {
	Name string
	URL  string
}

func (file *File) String() string {
	return fmt.Sprintf(
		"\n📎 [%s](%s)",
		lib.ParseTelegramMessage(file.Name),
		lib.ParseTelegramMessage(file.URL),
	)
}

// ParseGraphQL transforms a full-text GraphQL query into a json query containing it
func ParseGraphQL(query string) string {
	raw, _ := json.Marshal(map[string]string{
		"query": query,
	})
	return string(raw)
}
