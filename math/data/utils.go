package data

import (
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
)

// CheckExpression checks if a mathematical expression is valid (taking in account the syntax and available functions)
func CheckExpression(function string) (err error) {
	_, err = govaluate.NewEvaluableExpressionWithFunctions(function, Functions)
	if err != nil {
		err = fmt.Errorf("L'expression entréee est invalide : `%v`.", err)
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
