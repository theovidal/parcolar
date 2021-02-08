package api

import (
	"encoding/json"
	"fmt"
	"github.com/theovidal/bacbot/lib"
)

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

func ParseGraphql(query string) string {
	raw, _ := json.Marshal(map[string]string{
		"query": query,
	})
	return string(raw)
}
