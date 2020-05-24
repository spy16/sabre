package main

import (
	"fmt"
	"log"

	"github.com/spy16/sabre"
)

func main() {
	// Setup the environment available for your scripts. NewScope(nil)
	// starts with no bindings.
	scope := sabre.New()
	scope.Bind("api", sabre.ValueOf(&API{name: "foo"}))
	scope.Bind("console-print", sabre.ValueOf(printToConsole))
	scope.Bind("value-of-pi", sabre.ValueOf(valueOfPi))

	// userProgram can be read from a file, command-line, a network socket
	// etc. and can contain calls that return/simply have side effects.
	userProgram := `
		(api.SetName "Bob")
		(console-print (api.Name))
		(value-of-pi)
	`

	res, err := sabre.ReadEvalStr(scope, userProgram)
	if err != nil {
		panic(err)
	}

	fmt.Println(res) // should print 3.141200
}

func valueOfPi() float64 {
	return 3.1412
}

// You can expose control to your application through just functions
// also.
func printToConsole(msg string) {
	log.Printf("func-api called")
}

// API provides functions that allow your application behavior to be
// controlled at runtime.
type API struct {
	name string
}

// SetName can be used from the scripting layer to change name.
func (api *API) SetName(name string) {
	api.name = name
}

// Name returns the current value of the name.
func (api *API) Name() string {
	return api.name
}
