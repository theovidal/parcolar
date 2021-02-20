package lib

func Factorial(n float64) (o float64) {
	o = 1
	for i := 1.0; i <= n; i++ {
		o *= i
	}
	return
}
