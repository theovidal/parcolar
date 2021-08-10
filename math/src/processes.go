package src

import (
	"errors"
	"fmt"
	"math"

	"github.com/Knetic/govaluate"
)

// Integral uses Simpson's rule to integrate a function
// See https://en.wikipedia.org/wiki/Simpson%27s_rule
func Integral(args ...interface{}) (interface{}, error) {
	expression, ok := args[0].(string)
	if !ok {
		return nil, errors.New("L'expression doit être contenue dans une chaîne de caractères.")
	}
	if err := CheckExpression(expression); err != nil {
		return nil, err
	}

	a, ok := args[1].(float64)
	if !ok {
		return nil, nonRealError
	}
	b, ok := args[2].(float64)
	if !ok {
		return nil, nonRealError
	}
	n, ok := args[3].(float64)
	if !ok {
		return nil, nonRealError
	}

	var calculateSurface bool
	if len(args) > 4 {
		if surface, ok := args[4].(bool); ok {
			calculateSurface = surface
		}
	}

	h := (b - a) / n
	f, _ := govaluate.NewEvaluableExpressionWithFunctions(expression, BasicFunctions)
	var integral float64

	for i := 0.0; i <= n; i++ {
		x, m := a+h*i, a+h*(i+0.5)

		fX, err := evaluateForIntegral(f, calculateSurface, x)
		if err != nil {
			return nil, err
		}
		fM, err := evaluateForIntegral(f, calculateSurface, m)
		if err != nil {
			return nil, err
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

// evaluateForIntegral is a shortcut to calculate and verify a value specifically for an integral (real value wanted)
func evaluateForIntegral(f *govaluate.EvaluableExpression, surface bool, x float64) (y float64, err error) {
	variables := Constants
	variables["x"] = x
	eval, err := f.Evaluate(variables)
	y, ok := eval.(float64)
	if !ok {
		return y, nonRealReturnedError
	}

	if err != nil || y == math.NaN() || y == math.Inf(1) || y == math.Inf(-1) {
		return y, fmt.Errorf("Fonction non définie en %g.", x)
	}
	if y < 0 && surface {
		y *= -1
	}
	return
}
