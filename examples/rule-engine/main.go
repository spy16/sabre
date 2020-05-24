package main

import (
	"github.com/spy16/sabre"
	"github.com/spy16/sabre/core"
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
		// apply discount for the order
	} else {
		// don't apply discount
	}
}

func runDiscountingRule(rule string, user string) (bool, error) {
	// Define a scope with no bindings. (not even special forms)
	scope := sabre.New()

	bindGo := func(sym string, v interface{}) {
		_ = scope.Bind(sym, sabre.ValueOf(v))
	}

	// Define and expose your rules which ideally should have no
	// side effects.
	bindGo("and", and)
	bindGo("or", or)
	bindGo("regular-user?", isRegularUser)
	bindGo("minimum-cart-price?", isMinCartPrice)
	bindGo("not-blacklisted?", isNotBlacklisted)

	// Bind current user name
	bindGo("current-user", user)

	shouldDiscount, err := sabre.ReadEvalStr(scope, rule)
	return isTruthy(shouldDiscount), err
}

func isTruthy(v core.Value) bool {
	if v == nil || v == (core.Nil{}) {
		return false
	}
	if b, ok := v.(core.Bool); ok {
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
