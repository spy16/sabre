package sabre

// Module represents a group of forms. Evaluating a module leads to evaluation
// of each form in order and result will be the result of last evaluation.
type Module []Value

// Eval evaluates all the vals in the module body and returns the result of the
// last evaluation.
func (mod Module) Eval(scope Scope) (Value, error) {
	var res Value
	var err error

	for _, item := range mod {
		res, err = item.Eval(scope)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (mod Module) String() string {
	return containerString(mod, "", "\n", "\n")
}
