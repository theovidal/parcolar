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
	// Core curriculum
	"HISTOIRE & GEOGRAPHIE":             {Name: "Histoire-Géographie", Emoji: "🌎"},
	"ENS. MORAL & CIVIQUE":              {Name: "EMC", Emoji: "🏛"},
	"SC PHYSIQ & CHIMIQ":                {Name: "Sciences physiques", Emoji: "🔭"},
	"SCIENCES DE LA VIE ET DE LA TERRE": {Name: "SVT", Emoji: "☘"},
	"ED.PHYSIQUE & SPORT.":              {Name: "EPS", Emoji: "⚽"},
	"PHILOSOPHIE":                       {Name: "Philosophie", Emoji: "✒"},
	"PARCOURS REUSSITE ORIENT":          {Name: "MAP PRO (Vie de classe)", Emoji: "🪑"},

	// Living languages
	"ANGLAIS":            {Name: "Anglais", Emoji: "🍵"},
	"ESPAGNOL":           {Name: "Espagnol", Emoji: "🌮"},
	"DNL SI":             {Name: "Anglais Euro", Emoji: "🇪🇺"},
	"ANGLAIS SECT.EUROP": {Name: "Anglais Euro", Emoji: "🇪🇺"},

	// Specialties and options
	"MATHEMATIQUES":        {Name: "Mathématiques", Emoji: "🔢"},
	"MATHS EXP":            {Name: "Maths expertes", Emoji: "🧮"},
	"SC.INGEN. & SC.PHYS.": {Name: "Sciences de l'ingénieur", Emoji: "⚙"},
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
