package math

import (
	"errors"
	"github.com/Knetic/govaluate"
)

// CheckExpression checks if a mathematical expression is valid (taking in account the syntax and available functions)
func CheckExpression(function string) (err error) {
	_, err = govaluate.NewEvaluableExpressionWithFunctions(function, availableFunctions)
	if err != nil {
		err = errors.New("L'expression entréee est invalide : `" + err.Error() + "`.")
	}
	return
}

// Evaluate calculates f(x) for a certain function contained in an expression
func Evaluate(function string, x float64) (value float64, err error) {
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(function, availableFunctions)
	variables := availableConstants
	variables["x"] = x
	y, _ := expression.Evaluate(variables)
	value, ok := y.(float64)
	if !ok {
		err = errors.New("L'expression n'est pas valide, car certains symboles ne renvoient pas des nombres réels.")
	}
	return
}
