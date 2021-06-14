package data

import (
	"errors"
	"math"

	"github.com/Knetic/govaluate"
)

var Constants = map[string]interface{}{
	"e":   math.E,
	"pi":  math.Pi,
	"phi": math.Phi,
}

var BasicFunctions = map[string]govaluate.ExpressionFunction{
	// Classical
	"sqrt":  oneParamFunc(math.Sqrt),
	"abs":   oneParamFunc(math.Abs),
	"rem":   twoParamsFunc(math.Remainder),
	"gamma": oneParamFunc(math.Gamma),
	"fact":  oneParamFunc(Factorial),
	"floor": oneParamFunc(math.Floor),
	"ceil":  oneParamFunc(math.Ceil),

	// Logarithmic
	"exp": oneParamFunc(math.Exp),
	"ln":  oneParamFunc(math.Log),
	"log": oneParamFunc(math.Log10),

	// Trigonometry (radians)
	"sin":  oneParamFunc(math.Sin),
	"cos":  oneParamFunc(math.Cos),
	"tan":  oneParamFunc(math.Tan),
	"asin": oneParamFunc(math.Asin),
	"acos": oneParamFunc(math.Acos),
	"atan": oneParamFunc(math.Atan),
	"sinh": oneParamFunc(math.Sinh),
	"cosh": oneParamFunc(math.Cosh),
	"tanh": oneParamFunc(math.Tanh),

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
	"binomcoef": twoParamsFunc(BinomialCoefficient),
	"binomet":   threeParamsFunc(BinomialXEqualTo),
	"binomlt":   threeParamsFunc(BinomialXLessThan),
}

var Processes = map[string]govaluate.ExpressionFunction{
	// Calculus
	"integral": Integral,
	"surface":  Surface,
}

var Functions = map[string]govaluate.ExpressionFunction{}

// GatherFunctions lists the functions the user can use in their expressions
func GatherFunctions() {
	for name, process := range Processes {
		Functions[name] = process
	}
	for name, function := range BasicFunctions {
		Functions[name] = function
	}
}

// nonRealError is the default error shown when the user passes an argument to a numerical function that is not a real number
var nonRealError = errors.New("Merci de renseigner un nombre réel.")

// oneParamFunc is a short-hand helper to create a one-parameter mathematical function
func oneParamFunc(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x, ok := args[0].(float64)
		if !ok {
			return nil, nonRealError
		}
		return handler(x), nil
	}
}

// twoParamsFunc is a short-hand helper to create a two-parameters mathematical function
func twoParamsFunc(handler func(float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x, ok := args[0].(float64)
		if !ok {
			return nil, nonRealError
		}
		y, ok := args[1].(float64)
		if !ok {
			return nil, nonRealError
		}

		return handler(x, y), nil
	}
}

// threeParamsFunc is a short-hand helper to create a three-parameters mathematical function
func threeParamsFunc(handler func(float64, float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		a, ok := args[0].(float64)
		if !ok {
			return nil, nonRealError
		}
		b, ok := args[1].(float64)
		if !ok {
			return nil, nonRealError
		}
		c, ok := args[2].(float64)
		if !ok {
			return nil, nonRealError
		}

		return handler(a, b, c), nil
	}
}

// degreeToRadiansFunc is a short-hand helper to create a one-parameter mathematical function and convert the passed angle from degrees to radians (go native functions only admit this unit)
func degreeToRadiansFunc(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x, ok := args[0].(float64)
		if !ok {
			return nil, nonRealError
		}

		return handler(x * math.Pi / 180), nil
	}
}
