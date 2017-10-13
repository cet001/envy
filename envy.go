// Package envy is a little helper utility for interacting with the OS environment variables.
//
// A basic example is:
//
//  package main
//
//  import (
//      "fmt"
//      "github.com/cet001/envy"
//      "os"
//  )
//
//  func main() {
//      // Set some environment variables just for this example
//      os.Setenv("PORT", "9999")
//      os.Setenv("MY_TIMEOUT", "600")
//
//      env := envy.NewEnv()
//
//      port := env.Getenv("PORT").String()
//      timeout := env.Getenv("MY_TIMEOUT").Int()
//      foo := env.Getenv("FOO").DefaultValue("foo!")
//
//      // Check if any errors occured (e.g. int parsing error, attempting to get
//      // an environment variable that doesn't exist and no default value set)
//      if errors := env.Errors(); len(errors) > 0 {
//          for _, err := range errors {
//              fmt.Println(err)
//           }
//           os.Exit(1)
//      }
//
//      fmt.Printf("port=%v, timeout=%v, foo=%v\n", port, timeout, foo)
//  }
//
package envy

import (
	"fmt"
	"os"
	"strconv"
)

type Env struct {
	errors []error
}

func NewEnv() *Env {
	return &Env{
		errors: []error{},
	}
}

func (me *Env) Getenv(key string) Var {
	return Var{
		key:   key,
		value: os.Getenv(key),
		env:   me,
	}
}

// Returns any errors that occurred when getting the environment variables.
func (me *Env) Errors() []error {
	return me.errors
}

// Var represents an environment variable whose value can be cast to string or int.
type Var struct {
	key             string
	value           string
	hasDefaultValue bool
	env             *Env
}

// Determines the default value to use in this environment variable doesn't exist
// or is empty.
func (me Var) DefaultValue(defaultValue string) Var {
	if me.value == "" {
		me.value = defaultValue
	}
	me.hasDefaultValue = true
	return me
}

// Returns the string value of this environment variable.
func (me Var) String() string {
	if me.value == "" && !me.hasDefaultValue {
		me.env.errors = append(me.env.errors, fmt.Errorf("'%v' is empty, and no default value set.", me.key))
	}
	return me.value
}

// Returns the int value of this environment variable.
func (me Var) Int() int {
	value := me.String()
	if value == "" {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		me.env.errors = append(me.env.errors, fmt.Errorf("Error parsing '%v' as an int: %v", me.key, err))
	}
	return intValue
}
