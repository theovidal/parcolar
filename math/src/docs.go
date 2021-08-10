package src

import (
	"fmt"
	"strings"
)

// CalcDisclaimer is the default disclaimer for commands that use mathematical expressions
const CalcDisclaimer = "⚠ *Tous les signes multiplier* sont obligatoires (ex: 3x => 3 \\* x) et les *puissances* sont représentées par une *double-étoile* (\\*\\*).\nLes *fonctions trigonométriques* non précédées de la lettre `d` utilisent les *radians* comme unité pour les angles."

// DataDocumentation holds the documentation for the available functions and constants, to use in mathematical expressions
var DataDocumentation = func() string {
	var functionsDescription string
	for name := range BasicFunctions {
		functionsDescription += fmt.Sprintf("`%s`, ", name)
	}
	functionsDescription = strings.TrimSuffix(functionsDescription, ", ")

	var processesDescription string
	for name := range Processes {
		processesDescription += fmt.Sprintf("`%s`, ", name)
	}
	processesDescription = strings.TrimSuffix(processesDescription, ", ")

	var constantsDescription string
	for name := range Constants {
		constantsDescription += fmt.Sprintf("`%s`, ", name)
	}
	constantsDescription = strings.TrimSuffix(constantsDescription, ", ")

	return fmt.Sprintf("📈 Les fonctions disponibles sont : %s.\n⚙ Les procédés mathématiques, prenants en paramètres des expressions entre guillemets, sont : %s.\nπ Les constantes disponibles sont: %s.", functionsDescription, processesDescription, constantsDescription)
}()
