package data

import (
	"errors"
	"fmt"
	"math"

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
		x := a + h*i
		m := a + h*(i+0.5)

		variables["x"] = x
		eval, err := f.Evaluate(variables)
		fX := eval.(float64)
		if fX < 0 && calculateSurface {
			fX *= -1
		}

		variables["x"] = m
		eval, err = f.Evaluate(variables)
		fM := eval.(float64)
		if fM < 0 && calculateSurface {
			fM *= -1
		}

		if err != nil || fX == math.NaN() || fX == math.Inf(1) || fX == math.Inf(-1) {
			return nil, errors.New(fmt.Sprintf("Fonction non définie en %g", x))
		}

		if i < n {
			integral += 4 * fM
			if i > 0 {
				integral += 2 * fX
			}
		}
		if i == 0.0 || i == n {
			integral += fX
		}
	}

	return integral * h / 6, nil
}

// Surface is a short-hand to calculate integral of a function while taking in account its sign
func Surface(args ...interface{}) (interface{}, error) {
	args[4] = true
	return Integral(args)
}
