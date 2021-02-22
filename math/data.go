package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/Knetic/govaluate"

	"github.com/theovidal/bacbot/math/lib"
)

// mathFunc is a short-hand helper to create a one-parameter mathematical function
func mathFunc(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return handler(x), nil
	}
}

// mathFuncDouble is a short-hand helper to create a two-parameters mathematical function
func mathFuncDouble(handler func(float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		y := args[1].(float64)
		return handler(x, y), nil
	}
}

// mathFuncTriple is a short-hand helper to create a three-parameters mathematical function
func mathFuncTriple(handler func(float64, float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		a := args[0].(float64)
		b := args[1].(float64)
		c := args[2].(float64)
		return handler(a, b, c), nil
	}
}

// degreeToRadiansFunc is a short-hand helper to create a one-parameter mathematical function and convert the passed angle from degrees to radians (go native functions only admit this unit)
func degreeToRadiansFunc(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return handler(x * math.Pi / 180), nil
	}
}

// availableFunctions lists the functions the user can use in their expressions
var availableFunctions = map[string]govaluate.ExpressionFunction{
	// Classical
	"sqrt":  mathFunc(math.Sqrt),
	"abs":   mathFunc(math.Abs),
	"rem":   mathFuncDouble(math.Remainder),
	"gamma": mathFunc(math.Gamma),
	"fact":  mathFunc(lib.Factorial),

	// Logarithmic
	"exp": mathFunc(math.Exp),
	"ln":  mathFunc(math.Log),
	"log": mathFunc(math.Log10),

	// Trigonometry (radians)
	"sin":  mathFunc(math.Sin),
	"cos":  mathFunc(math.Cos),
	"tan":  mathFunc(math.Tan),
	"asin": mathFunc(math.Asin),
	"acos": mathFunc(math.Acos),
	"atan": mathFunc(math.Atan),
	"sinh": mathFunc(math.Sinh),
	"cosh": mathFunc(math.Cosh),
	"tanh": mathFunc(math.Tanh),

	// Trigonometry (degrees)
	"dsin":  degreeToRadiansFunc(math.Sin),
	"dcos":  degreeToRadiansFunc(math.Cos),
	"dtan":  degreeToRadiansFunc(math.Tan),
	"dasin": degreeToRadiansFunc(math.Asin),
	"dacos": degreeToRadiansFunc(math.Acos),
	"datan": degreeToRadiansFunc(math.Atan),
	"dsinh": degreeToRadiansFunc(math.Sinh),
	"dcosh": degreeToRadiansFunc(math.Cosh),
	"dtanh": degreeToRadiansFunc(math.Tanh),

	// Probabilities
	"binomcoef": mathFuncDouble(lib.BinomialCoefficient),
	"binomet":   mathFuncTriple(lib.BinomialXEqualTo),
	"binomlt":   mathFuncTriple(lib.BinomialXLessThan),
}

// availableConstants lists the constants the user can use in their expressions
var availableConstants = map[string]interface{}{
	"e":   math.E,
	"pi":  math.Pi,
	"phi": math.Phi,
}

// calcDisclaimer is the default disclaimer for commands that use mathematical expressions
const calcDisclaimer = "⚠ *Tous les signes multiplier* sont obligatoires (ex: 3x => 3 \\* x) et les *puissances* sont représentées par une *double-étoile* (\\*\\*).\nLes *fonctions trigonométriques non précédées de la lettre `d` utilisent les *radians* comme unité pour les angles."

// dataDocumentation holds the documentation for the available functions and constants, to use in mathematical expressions
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
