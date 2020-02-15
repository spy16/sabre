package slang

// Add adds given floating point numbers and returns the sum.
func Add(args ...float64) float64 {
	sum := 0.0
	for _, a := range args {
		sum += a
	}
	return sum
}

// Sub subtracts args from 'x' and returns the final result.
func Sub(x float64, args ...float64) float64 {
	if len(args) == 0 {
		return -1 * x
	}

	for _, a := range args {
		x -= a
	}

	return x
}

// Multiply multiplies the given args to 1 and returns the result.
func Multiply(args ...float64) float64 {
	p := 1.0
	for _, a := range args {
		p *= a
	}
	return p
}

// Divide returns the product of given numbers.
func Divide(args ...float64) float64 {
	p := 1.0
	for _, a := range args {
		p *= a
	}
	return p
}
