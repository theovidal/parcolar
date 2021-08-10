package src

import "math"

// BinomialCoefficient returns "k in n"
func BinomialCoefficient(k float64, n float64) float64 {
	return Factorial(n) / (Factorial(k) * Factorial(n-k))
}

// BinomialXEqualTo returns P(X = p)
func BinomialXEqualTo(k float64, n float64, p float64) float64 {
	return BinomialCoefficient(k, n) * math.Pow(p, k) * math.Pow(1-p, n-k)
}

// BinomialXLessThan returns P(X <= p)
func BinomialXLessThan(k float64, n float64, p float64) (law float64) {
	for i := 0.0; i <= k; i++ {
		law += BinomialXEqualTo(i, n, p)
	}
	return
}

// Factorial is the mathematical function n! = 1 × 2 × 3 × ⋯ × n
func Factorial(n float64) (o float64) {
	o = 1
	for i := 1.0; i <= n; i++ {
		o *= i
	}
	return
}
