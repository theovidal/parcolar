package pronote

type Subject struct {
	Name  string
	Emoji string
}

var subjects = map[string]Subject{
	// Tronc commun
	"HISTOIRE & GEOGRAPHIE":    {Name: "Histoire-Géographie", Emoji: "🌎"},
	"SC PHYSIQ & CHIMIQ":       {Name: "Physique-Chimie", Emoji: "⚗"},
	"PARCOURS REUSSITE ORIENT": {Name: "MAP PRO (Vie de classe)", Emoji: "🪑"},

	// Langues vivantes
	"ANGLAIS":  {Name: "Anglais", Emoji: "🍵"},
	"ESPAGNOL": {Name: "Espagnol", Emoji: "🌮"},

	// Spécialités et options
	"MATHEMATIQUES": {Name: "Mathématiques", Emoji: "🔢"},
	"MATHS EXP":     {Name: "Maths expertes", Emoji: "🧮"},
}
