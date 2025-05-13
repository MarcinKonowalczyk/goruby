package interpreter

import (
	"fmt"
	"go/token"
	"os"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/trace"
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
	node, tracer, err := parser.ParseFileEx(token.NewFileSet(), filename, input)
	if i.Trace {
		// if len(tracer_messages) > 0 {
		// 	fmt.Println("Tracer messages:")
		// 	for _, message := range tracer_messages {
		// 		fmt.Println(message)
		// 	}
		// 	fmt.Println("End of tracer messages")
		// }
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		indent := 0
		walkable.Walk(func(n trace.Node) error {
			if n.Name() == trace.START_NODE || n.Name() == trace.END_NODE {
				// ignore start and end nodes
				return nil
			}
			switch n.(type) {
			case *trace.Enter:
				fmt.Printf("%s > %s\n", strings.Repeat(".", indent*2), n.Name())
				indent++
			case *trace.Exit:
				indent--
				fmt.Printf("%s < %s\n", strings.Repeat(".", indent*2), n.Name())
				//
			}
			return nil
		})
	}
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return evaluator.Eval(node, i.environment)
}
