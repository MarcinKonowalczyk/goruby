package interpreter

import (
	"go/token"
	"log"
	"os"

	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/parser"
)

// // Interpreter defines the methods of an interpreter
// type Interpreter interface {
// 	Interpret(filename string, input interface{}) (object.RubyObject, error)
// }

// NewInterpreter returns an Interpreter ready to use and with the environment set to
// object.NewMainEnvironment()
func NewInterpreter() Interpreter {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get working directory: %s\n", err)
	}
	env := object.NewMainEnvironment()
	loadPath, _ := env.Get("$:")
	loadPathArr := loadPath.(*object.Array)
	loadPathArr.Elements = append(loadPathArr.Elements, &object.String{Value: cwd})
	env.SetGlobal("$:", loadPathArr)

	// setup ARGV
	argv := os.Args[1:]
	if len(argv) > 0 {
		// take off one more element
		argv = argv[1:]
	}
	argvArr := object.NewArray()
	for _, arg := range argv {
		argvArr.Elements = append(argvArr.Elements, &object.String{Value: arg})
	}
	env.SetGlobal("ARGV", argvArr)
	return Interpreter{environment: env}
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
