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
func Divide(first float64, args ...float64) float64 {
	if len(args) == 0 {
		return 1 / first
	}

	for _, a := range args {
		first /= a
	}

	return first
}

// Lt returns true if the given args are monotonically increasing.
func Lt(base float64, args ...float64) bool {
	inc := true
	for _, arg := range args {
		inc = inc && (arg > base)
	}
	return inc
}

// LtE returns true if the given args are monotonically increasing or
// are all equal.
func LtE(base float64, args ...float64) bool {
	inc := true
	for _, arg := range args {
		inc = inc && (arg >= base)
	}
	return inc
}

// Gt returns true if the given args are monotonically decreasing.
func Gt(base float64, args ...float64) bool {
	inc := true
	for _, arg := range args {
		inc = inc && (arg < base)
	}
	return inc
}

// GtE returns true if the given args are monotonically decreasing or
// all equal.
func GtE(base float64, args ...float64) bool {
	inc := true
	for _, arg := range args {
		inc = inc && (arg <= base)
	}
	return inc
}
