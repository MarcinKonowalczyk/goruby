package interpreter

import (
	"os"

	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/parser"
	parser_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/parser"
	"github.com/MarcinKonowalczyk/goruby/trace/printer"
	"github.com/MarcinKonowalczyk/goruby/transformer"
)

// // Interpreter defines the methods of an interpreter
type Interpreter interface {
	Interpret(filename string) (object.RubyObject, error)
	InterpretCode(code string) (object.RubyObject, error)
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
	environment     object.Environment
	trace_parse     bool
	trace_transform bool
	trace_eval      bool
}

func (i *interpreter) InterpretCode(src string) (object.RubyObject, error) {
	program, err := parser.Parse(src)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}

	// const ENABLE_TRANSFORMS_IN_INTERPRETER = true
	const ENABLE_TRANSFORMS_IN_INTERPRETER = false
	if ENABLE_TRANSFORMS_IN_INTERPRETER {
		program, err = transformer.Transform(program, transformer.ALL_STAGES, i.trace_transform)
		if err != nil {
			return nil, object.NewRuntimeError("transformer error: %v", err)
		}
	}

	res, tracer, err := evaluator.EvalEx(program, i.environment, i.trace_eval)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(printer.NewTracePrinter())
	}
	return res, err
}

func (i *interpreter) Interpret(filename string) (object.RubyObject, error) {
	program, tracer, err := parser_pipeline.ParseFile(filename, i.trace_parse)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(printer.NewTracePrinter())
	}
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}

	// const ENABLE_TRANSFORMS_IN_INTERPRETER = true
	const ENABLE_TRANSFORMS_IN_INTERPRETER = false
	if ENABLE_TRANSFORMS_IN_INTERPRETER {
		program, err = transformer.Transform(program, transformer.ALL_STAGES, i.trace_transform)
		if err != nil {
			return nil, object.NewRuntimeError("transformer error: %v", err)
		}
	}

	res, tracer, err := evaluator.EvalEx(program, i.environment, i.trace_eval)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(printer.NewTracePrinter())
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
