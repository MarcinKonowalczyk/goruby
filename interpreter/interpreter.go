package interpreter

import (
	"go/token"
	"os"

	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/parser"
)

// // Interpreter defines the methods of an interpreter
// type Interpreter interface {
// 	Interpret(filename string, input interface{}) (object.RubyObject, error)
// }

func NewInterpreterEx(argv []string) Interpreter {
	env := object.NewMainEnvironment()

	argvArr := object.NewArray()
	for _, arg := range argv {
		argvArr.Elements = append(argvArr.Elements, object.NewString(arg))
	}
	env.SetGlobal("ARGV", argvArr)
	return Interpreter{environment: env}
}

// NewInterpreter returns an Interpreter ready to use and with the environment set to
// object.NewMainEnvironment()
func NewInterpreter() Interpreter {
	return NewInterpreterEx(os.Args[1:])
}

type Interpreter struct {
	environment object.Environment
	Trace       bool
}

func (i *Interpreter) Interpret(filename string, input interface{}) (object.RubyObject, error) {
	var mode parser.Mode
	if i.Trace {
		mode = parser.Trace
	}
	node, err := parser.ParseFile(token.NewFileSet(), filename, input, mode)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return evaluator.Eval(node, i.environment)
}
