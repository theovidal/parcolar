package math

import (
	"fmt"
	"github.com/Knetic/govaluate"
)

// CheckExpression checks if a mathematical expression is valid (taking in account the syntax and available functions)
func CheckExpression(function string) (message string) {
	fmt.Println(function)
	if len(function) == 0 {
		message = "Merci d'indiquer une expression mathématique."
		return
	}
	_, err := govaluate.NewEvaluableExpressionWithFunctions(function, availableFunctions)
	if err != nil {
		message = "L'expression entréee est invalide : `" + err.Error() + "`."
	}
	return
}

// Evaluate calculates f(x) for a certain function contained in an expression
func Evaluate(function string, x float64) float64 {
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(function, availableFunctions)
	variables := availableConstants
	variables["x"] = x
	y, _ := expression.Evaluate(variables)
	return y.(float64)
}
