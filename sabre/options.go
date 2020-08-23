package sabre

// Option can be used with New() to customize initialization of Sabre.
type Option func(s *Sabre)

// WithRuntime sets a custom Runtime value while initializing Sabre. If rt
// is nil, uses the default Runtime value.
func WithRuntime(rt Runtime) Option {
	return func(s *Sabre) {
		if rt != nil {
			s.rt = rt
		}
	}
}

// WithMaxDepth sets the maximum stack depth allowed for function calls in
// Sabre instance.
func WithMaxDepth(depth int) Option {
	return func(s *Sabre) {
		s.maxDepth = depth
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithRuntime(nil),
	}, opts...)
}
