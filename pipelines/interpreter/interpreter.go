package interpreter

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/evaluator"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/MarcinKonowalczyk/goruby/trace/printer"
	"github.com/MarcinKonowalczyk/goruby/transformer"
)

// // Interpreter defines the methods of an interpreter
type Interpreter interface {
	Interpret(filename string) (ruby.Object, error)
	InterpretCode(code string) (ruby.Object, error)
	ParseCode(src string) (*ast.Program, error)
	// SetTraceParse sets the trace_parse flag
	SetTraceParse(trace_parse bool)
	// SetTraceEval sets the trace_eval flag
	SetTraceEval(trace_eval bool)
}

func NewInterpreter(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	argv []string) Interpreter {
	env := object.NewMainEnvironment()

	argvArr := object.NewArray()
	for _, arg := range argv {
		argvArr.Elements = append(argvArr.Elements, object.NewString(arg))
	}
	env.SetGlobal("ARGV", argvArr)
	return &interpreter{
		stdin:       stdin,
		stdout:      stdout,
		stderr:      stderr,
		environment: env,
	}
}

func NewBasicInterpreter() Interpreter {
	return NewInterpreter(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Args[1:],
	)
}

type interpreter struct {
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	environment env.Environment[ruby.Object]

	// trace parsing flags
	trace_parse          bool
	print_parse_messages bool

	trace_transform bool
	trace_eval      bool
}

func (i *interpreter) ParseCode(src string) (*ast.Program, error) {
	ctx := context.Background()
	var parse_tracer trace.Tracer
	if i.trace_parse {
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
		var out strings.Builder
		_ = walkable.Walk(printer.NewTracePrinter(&out, i.print_parse_messages))
		i.stdout.Write([]byte(out.String()))
	}

	return program, err
}

func (i *interpreter) InterpretCode(src string) (ruby.Object, error) {
	program, err := i.ParseCode(src)

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

	ctx := context.Background()
	var eval_tracer trace.Tracer
	if i.trace_eval {
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
		_ = walkable.Walk(printer.NewTracePrinter(os.Stdout, true))
	}

	return res, err
}

func (i *interpreter) Interpret(filename string) (ruby.Object, error) {

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, object.NewRuntimeError("could not read file %s: %v", filename, err)
	}

	program, err := i.ParseCode(string(src))

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

	ctx := context.Background()
	var eval_tracer trace.Tracer
	if i.trace_eval {
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
		_ = walkable.Walk(printer.NewTracePrinter(os.Stdout, true))
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
