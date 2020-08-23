package sabre

// ParseSpecial validates a special form invocation, parse the form and
// returns an expression that can be evaluated for result.
type ParseSpecial func(s *Sabre, args Seq) (Expr, error)
