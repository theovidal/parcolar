package lib

import "math"

func BinomialCoefficient(k float64, n float64) float64 {
	return Factorial(n) / (Factorial(k) * Factorial(n-k))
}

func BinomialXEqualTo(k float64, n float64, p float64) float64 {
	return BinomialCoefficient(k, n) * math.Pow(p, k) * math.Pow(1-p, n-k)
}

func BinomialXLessThan(k float64, n float64, p float64) (law float64) {
	for i := 0.0; i <= k; i++ {
		law += BinomialXEqualTo(i, n, p)
	}
	return
}
