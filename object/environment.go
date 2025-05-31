package object

import (
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
)

var (
	CLASSES = env.NewEnvironment[ruby.Object]()
	// symbol to which all function definitions are attached to
	FUNCS_STORE = newExtendedObject(FUNCS)
)

// NewMainEnvironment returns a new Environment populated with all Ruby classes
// and the Kernel functions
func NewMainEnvironment() env.Environment[ruby.Object] {
	env := CLASSES.Clone()
	env.Set("bottom", BOTTOM)
	env.Set("funcs", FUNCS_STORE)
	env.SetGlobal("$stdin", IoClass)
	return env
}
