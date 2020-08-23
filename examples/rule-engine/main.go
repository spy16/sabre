package main

import (
	"fmt"
	"strings"

	"github.com/spy16/sabre"
	"github.com/spy16/sabre/runtime"
)

func main() {
	// Accept business rules from file, command-line, http request etc.
	// These rules can change as per business requirements and your
	// application doesn't have to change.
	ruleSrc := `(and (regular-user? current-user)
					 (not-blacklisted? current-user))`

	shouldDiscount, err := runDiscountingRule(ruleSrc, "bob")
	if err != nil {
		panic(err)
	}

	if shouldDiscount {
		fmt.Println("apply discount")
	} else {
		fmt.Println("don't apply discount")
	}
}

func runDiscountingRule(rule string, user string) (bool, error) {
	// Define a runtime with no bindings. (not even special forms)
	rt := sabre.New()

	bind := func(sym string, v interface{}) {
		_ = rt.Bind(sym, sabre.ValueOf(v))
	}

	// Define and expose your rules which ideally should have no
	// side effects.
	bind("and", and)
	bind("regular-user?", isRegularUser)
	bind("minimum-cart-price?", isMinCartPrice)
	bind("not-blacklisted?", isNotBlacklisted)

	// bind current user name
	bind("current-user", user)

	shouldDiscount, err := sabre.ReadEval(rt, strings.NewReader(rule))
	return isTruthy(shouldDiscount), err
}

func isTruthy(v runtime.Value) bool {
	if v == nil || v == (runtime.Nil{}) {
		return false
	}
	if b, ok := v.(runtime.Bool); ok {
		return bool(b)
	}
	return true
}

func isNotBlacklisted(user string) bool {
	return user != "joe"
}

func isMinCartPrice(price float64) bool {
	return price >= 100
}

func isRegularUser(user string) bool {
	return user == "bob"
}

func and(rest ...bool) bool {
	if len(rest) == 0 {
		return true
	}
	result := rest[0]
	for _, r := range rest {
		result = result && r
		if !result {
			return false
		}
	}
	return true
}

func or(rest ...bool) bool {
	if len(rest) == 0 {
		return true
	}
	result := rest[0]
	for _, r := range rest {
		if result {
			return true
		}
		result = result || r
	}
	return false
}
