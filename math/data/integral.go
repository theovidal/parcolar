package data

import (
	"errors"
	"fmt"

	"github.com/Knetic/govaluate"
)

func Integral(args ...interface{}) (interface{}, error) {
	expression, ok := args[0].(string)
	if !ok {
		return nil, errors.New("L'expression doit être contenue dans une chaîne de caractères.")
	}
	if err := CheckExpression(expression); err != nil {
		return nil, err
	}

	a := args[1].(float64)
	b := args[2].(float64)
	n := args[3].(float64)

	var calculateSurface bool
	if len(args) > 4 {
		surface, ok := args[4].(bool)
		if ok {
			calculateSurface = surface
		}
	}

	h := (b - a) / n

	f, _ := govaluate.NewEvaluableExpressionWithFunctions(expression, BasicFunctions)
	variables := Constants

	var integral float64
	for i := 0.0; i <= n; i++ {
		variables["x"] = a + h*i
		eval, err := f.Evaluate(variables)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Fonction non définie en %s", eval))
		}

		value := eval.(float64)
		if i == 0.0 || i == n {
			value *= 0.5
		}
		if value < 0 && calculateSurface {
			value *= -1
		}
		integral += value
	}

	return integral * h, nil
}

func Surface(args ...interface{}) (interface{}, error) {
	args[4] = true
	return Integral(args)
}
