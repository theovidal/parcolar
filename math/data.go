package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/Knetic/govaluate"
)

func mathFunc(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return handler(x), nil
	}
}

func mathFuncDouble(handler func(float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		y := args[1].(float64)
		return handler(x, y), nil
	}
}

var availableFunctions = map[string]govaluate.ExpressionFunction{
	"sqrt":  mathFunc(math.Sqrt),
	"abs":   mathFunc(math.Abs),
	"rem":   mathFuncDouble(math.Remainder),
	"gamma": mathFunc(math.Gamma),

	"exp": mathFunc(math.Exp),
	"ln":  mathFunc(math.Log),
	"log": mathFunc(math.Log10),

	"sin":  mathFunc(math.Sin),
	"cos":  mathFunc(math.Cos),
	"tan":  mathFunc(math.Tan),
	"asin": mathFunc(math.Asin),
	"acos": mathFunc(math.Acos),
	"atan": mathFunc(math.Atan),
	"sinh": mathFunc(math.Sinh),
	"cosh": mathFunc(math.Cosh),
	"tanh": mathFunc(math.Tanh),
}

var availableConstants = map[string]interface{}{
	"e":   math.E,
	"pi":  math.Pi,
	"phi": math.Phi,
}

const calcDisclaimer = "⚠ *Tous les signes multiplier* sont obligatoires (ex: 3x => 3 \\* x) et les *puissances* sont représentées par une *double-étoile* (\\*\\*).\nLes *fonctions trigonométriques* utilisent toutes les *radians* comme unité pour les angles."

var dataDocumentation = func() string {
	var functionsDescription string
	for name := range availableFunctions {
		functionsDescription += fmt.Sprintf("`%s`, ", name)
	}
	functionsDescription = strings.TrimSuffix(functionsDescription, ", ")

	var constantsDescription string
	for name := range availableConstants {
		constantsDescription += fmt.Sprintf("`%s`, ", name)
	}
	constantsDescription = strings.TrimSuffix(constantsDescription, ", ")

	return fmt.Sprintf("📈 Les fonctions disponibles sont : %s.\nπ Les constantes disponibles sont: %s.", functionsDescription, constantsDescription)
}()
