package repl

import (
	log "github.com/lthibault/log/pkg"
)

// Option implmentations can be provided to New to configure the
// REPL during initialization.
type Option func(repl *REPL)

// WithLogger sets the REPL's logger.  `nil` is a no-op logger.
func WithLogger(log log.Logger) Option {
	return func(repl *REPL) {
		repl.log = log
	}
}

// WithPrompt sets the REPL's prompt.  `nil` uses a libreadline implementation.
func WithPrompt(prompt Prompt) Option {
	return func(repl *REPL) {
		repl.prompt = prompt
	}
}

// WithBanner sets the REPL's banner.
func WithBanner(banner string) Option {
	return func(repl *REPL) {
		repl.banner = banner
	}
}

func withDefaults(opt []Option) []Option {
	return append([]Option{
		WithLogger(nil),
		WithBanner("SLANG - a tiny lisp based on Sabre."),
		// WithSomeOtherOption(...)
	}, opt...)
}

// Prompt signals that a goroutine is ready to accept input by setting a prompt, and
// reads in a line, blocking until one is available.
type Prompt interface {
	SetPrompt(string)
	Readline() (string, error)
}
