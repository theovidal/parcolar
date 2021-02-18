package math

import (
	"math"

	"github.com/Knetic/govaluate"
)

func function(handler func(x float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return handler(x), nil
	}
}

func doubleFunction(handler func(float64, float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		y := args[1].(float64)
		return handler(x, y), nil
	}
}

var functions = map[string]govaluate.ExpressionFunction{
	"sqrt":  function(math.Sqrt),
	"abs":   function(math.Abs),
	"rem":   doubleFunction(math.Remainder),
	"gamma": function(math.Gamma),

	"exp": function(math.Exp),
	"ln":  function(math.Log),
	"log": function(math.Log10),

	"sin":  function(math.Sin),
	"cos":  function(math.Cos),
	"tan":  function(math.Tan),
	"asin": function(math.Asin),
	"acos": function(math.Acos),
	"atan": function(math.Atan),
	"sinh": function(math.Sinh),
	"cosh": function(math.Cosh),
	"tanh": function(math.Tanh),
}

var constants = map[string]interface{}{
	"e":   math.E,
	"pi":  math.Pi,
	"phi": math.Phi,
}
