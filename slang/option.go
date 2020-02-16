package slang

import (
	log "github.com/lthibault/log/pkg"
)

// REPLOption implmentations can be provided to NewREPL to configure the
// REPL during initialization.
type REPLOption func(repl *REPL)

// WithLogger sets the REPL's logger.  `nil` is a no-op logger.
func WithLogger(log log.Logger) REPLOption {
	return func(repl *REPL) {
		repl.log = log
	}
}

// WithPrompt sets the REPL's prompt.  `nil` uses a libreadline implementation.
func WithPrompt(prompt Prompt) REPLOption {
	return func(repl *REPL) {
		repl.prompt = prompt
	}
}

// WithBanner sets the REPL's banner.
func WithBanner(banner string) REPLOption {
	return func(repl *REPL) {
		repl.banner = banner
	}
}

func withDefaults(opt []REPLOption) []REPLOption {
	return append([]REPLOption{
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

// RuntimeOption implementations can be provided to New() to configure
// the language runtime.
type RuntimeOption func(slang *Slang)

// RuntimeLogger sets the runtime's logger instance.
func RuntimeLogger(log log.Logger) RuntimeOption {
	return func(slang *Slang) {
		slang.log = log
	}
}
