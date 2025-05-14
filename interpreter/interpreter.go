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
type Interpreter interface {
	Interpret(filename string, input interface{}) (object.RubyObject, error)
	// SetTraceParse sets the trace_parse flag
	SetTraceParse(trace_parse bool)
	// SetTraceEval sets the trace_eval flag
	SetTraceEval(trace_eval bool)
}

func NewInterpreterEx(argv []string) Interpreter {
	env := object.NewMainEnvironment()

	argvArr := object.NewArray()
	for _, arg := range argv {
		argvArr.Elements = append(argvArr.Elements, object.NewString(arg))
	}
	env.SetGlobal("ARGV", argvArr)
	return &interpreter{environment: env}
}

// NewInterpreter returns an Interpreter ready to use and with the environment set to
// object.NewMainEnvironment()
func NewInterpreter() Interpreter {
	return NewInterpreterEx(os.Args[1:])
}

type interpreter struct {
	environment object.Environment
	trace_parse bool
	trace_eval  bool
}

func newTracePrinter() func(trace.Node) error {
	indent := 0
	return func(n trace.Node) error {
		switch n := n.(type) {
		case *trace.Enter:
			if n.Name == trace.START_NODE {
				return nil
			}
			fmt.Printf("%s > %s\n", strings.Repeat(".", indent*2), n.Name)
			indent++
		case *trace.Exit:
			if n.Name == trace.END_NODE {
				return nil
			}
			indent--
			fmt.Printf("%s < %s\n", strings.Repeat(".", indent*2), n.Name)
		case *trace.Message:
			fmt.Printf("%s %s\n", strings.Repeat(".", indent*2), n.Message)
		default:
			panic(fmt.Sprintf("unknown node type: %T", n))
		}
		return nil
	}
}

func (i *interpreter) Interpret(filename string, input interface{}) (object.RubyObject, error) {
	node, tracer, err := parser.ParseFileEx(token.NewFileSet(), filename, input, i.trace_parse)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(newTracePrinter())
	}
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	res, tracer, err := evaluator.EvalEx(node, i.environment, i.trace_eval)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(newTracePrinter())
	}
	return res, err
}

func (i *interpreter) SetTraceParse(trace_parse bool) {
	i.trace_parse = trace_parse
}
func (i *interpreter) SetTraceEval(trace_eval bool) {
	i.trace_eval = trace_eval
}

var _ Interpreter = &interpreter{}
