package math

import (
	"math"

	"github.com/Knetic/govaluate"
)

var functions = map[string]govaluate.ExpressionFunction{
	"sqrt": func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return math.Sqrt(x), nil
	},
	"exp": func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return math.Exp(x), nil
	},
	"ln": func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return math.Log(x), nil
	},
	"sin": func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return math.Sin(x), nil
	},
	"cos": func(args ...interface{}) (interface{}, error) {
		x := args[0].(float64)
		return math.Cos(x), nil
	},
}
