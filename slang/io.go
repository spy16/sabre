package slang

import (
	"fmt"
)

// Println is an alias for fmt.Println which ignores the return values.
func Println(args ...interface{}) error {
	_, err := fmt.Println(args...)
	return err
}

// Printf is an alias for fmt.Printf which ignores the return values.
func Printf(format string, args ...interface{}) error {
	_, err := fmt.Printf(format, args...)
	return err
}
