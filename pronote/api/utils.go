package api

import (
	"encoding/json"
	"fmt"

	"github.com/theovidal/parcolar/lib"
)

// Subject defines the pretty-print style of a school subject
type Subject struct {
	Name  string
	Emoji string
}

// Subjects holds the list of all available subjects in Pronote
var Subjects = map[string]Subject{
	"MATHÉMATIQUES":        {Name: "Mathématiques", Emoji: "🔢"},
	"PHYSIQUE":                {Name: "Sciences physiques", Emoji: "🔭"},
	"INFO COURS":     {Name: "Informatique", Emoji: "💻"},
	"INFO TP":        {Name: "Informatique (TP)", Emoji: "🖥"},

	"FRANÇAIS PHILO": {Name: "Français/Philo", Emoji: "✒"},
	"ANGLAIS":            {Name: "Anglais", Emoji: "🍵"},
	"ESPAGNOL":           {Name: "Espagnol", Emoji: "🌮"},

	"SCIENCES INGÉNIEUR": {Name: "Sciences de l'ingénieur", Emoji: "⚙"},
	"SC ING TD": {Name: "Sciences de l'ingénieur (TD)", Emoji: "⚙"},

	"DEVOIRS": {Name: "Devoir surveillé", Emoji: "✏"},
}

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
