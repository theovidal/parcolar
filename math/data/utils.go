package data

import (
	"errors"
	"math"

	"github.com/Knetic/govaluate"
)

// CheckExpression checks if a mathematical expression is valid (taking in account the syntax and available functions)
func CheckExpression(function string) (err error) {
	_, err = govaluate.NewEvaluableExpressionWithFunctions(function, Functions)
	if err != nil {
		err = errors.New("L'expression entréee est invalide : `" + err.Error() + "`.")
	}
	return
}

// Evaluate calculates f(x) for a certain function contained in an expression
func Evaluate(function string, x float64) (value float64, err error) {
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(function, Functions)
	variables := Constants
	variables["x"] = x
	y, err := expression.Evaluate(variables)
	if err != nil {
		return 0.0, err
	}
	value, ok := y.(float64)
	if !ok {
		err = errors.New("L'expression n'est pas valide, car certains symboles ne renvoient pas des nombres réels.")
	}
	return
}

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
