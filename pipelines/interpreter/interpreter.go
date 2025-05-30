package interpreter

import (
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/parser"
	parser_pipeline "github.com/MarcinKonowalczyk/goruby/pipelines/parser"
	"github.com/MarcinKonowalczyk/goruby/trace"
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

func NewInterpreter(argv []string) Interpreter {
	env := object.NewMainEnvironment()

	argvArr := object.NewArray()
	for _, arg := range argv {
		argvArr.Elements = append(argvArr.Elements, object.NewString(arg))
	}
	env.SetGlobal("ARGV", argvArr)
	return &interpreter{environment: env}
}

type interpreter struct {
	environment     object.Environment
	trace_parse     bool
	trace_transform bool
	trace_eval      bool
}

func (i *interpreter) InterpretCode(src string) (object.RubyObject, error) {
	ctx := context.Background()
	var parse_tracer trace.Tracer
	if i.trace_parse {
		fmt.Printf("tracing")
		parse_tracer = trace.NewTracer()
		ctx = trace.WithTracer(ctx, parse_tracer)
	}
	program, err := parser.ParseCtx(ctx, src)
	if parse_tracer != nil {
		parse_tracer.Done()
		walkable, err := parse_tracer.ToWalkable()
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

	var eval_tracer trace.Tracer
	if i.trace_eval {
		fmt.Printf("tracing")
		eval_tracer = trace.NewTracer()
		ctx = trace.WithTracer(ctx, eval_tracer)
	}
	res, err := evaluator.Eval(ctx, program, i.environment)
	if eval_tracer != nil {
		eval_tracer.Done()
		walkable, err := eval_tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(printer.NewTracePrinter())
	}

	return res, err
}

func (i *interpreter) Interpret(filename string) (object.RubyObject, error) {
	ctx := context.Background()
	var parse_tracer trace.Tracer
	if i.trace_parse {
		parse_tracer = trace.NewTracer()
		ctx = trace.WithTracer(ctx, parse_tracer)
	}
	program, err := parser_pipeline.ParseFile(ctx, filename)
	if parse_tracer != nil {
		parse_tracer.Done()
		walkable, err := parse_tracer.ToWalkable()
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

	var eval_tracer trace.Tracer
	if i.trace_eval {
		fmt.Printf("tracing")
		eval_tracer = trace.NewTracer()
		ctx = trace.WithTracer(ctx, eval_tracer)
	}
	res, err := evaluator.Eval(ctx, program, i.environment)
	if eval_tracer != nil {
		eval_tracer.Done()
		walkable, err := eval_tracer.ToWalkable()
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
