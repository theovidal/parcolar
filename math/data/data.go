package data

import (
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
	"sqrt":  mathFunc(math.Sqrt),
	"abs":   mathFunc(math.Abs),
	"rem":   mathFuncDouble(math.Remainder),
	"gamma": mathFunc(math.Gamma),
	"fact":  mathFunc(Factorial),

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
	"binomcoef": mathFuncDouble(BinomialCoefficient),
	"binomet":   mathFuncTriple(BinomialXEqualTo),
	"binomlt":   mathFuncTriple(BinomialXLessThan),
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
