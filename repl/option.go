package repl

// Option implmentations can be provided to New to configure the
// REPL during initialization.
type Option func(repl *REPL)

// WithInput sets the REPL's input.  `nil` uses a libreadline implementation.
func WithInput(input Inputter) Option {
	return func(repl *REPL) {
		repl.input = input
	}
}

// WithBanner sets the REPL's banner.
func WithBanner(banner string) Option {
	return func(repl *REPL) {
		repl.banner = banner
	}
}

// WithPrompts sets the prompt to be displayed when waiting for user input
// in the REPL.
func WithPrompts(oneLine, multiLine string) Option {
	return func(repl *REPL) {
		repl.prompt = oneLine
		repl.multiPrompt = multiLine
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithBanner("SLANG - a tiny lisp based on Sabre."),
		WithPrompts("=>", "|"),
		// WithSomeOtherOption(...)
	}, opt...)
}

// Inputter signals that a goroutine is ready to accept input by setting a prompt, and
// reads in a line, blocking until one is available.
type Inputter interface {
	SetPrompt(string)
	Readline() (string, error)
}
