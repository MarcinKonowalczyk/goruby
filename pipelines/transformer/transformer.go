package transformer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"

	node_printer "github.com/MarcinKonowalczyk/goruby/printer"
	t "github.com/MarcinKonowalczyk/goruby/transformer"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
	"github.com/MarcinKonowalczyk/trace"
	"github.com/MarcinKonowalczyk/trace/printer"
)

type Transformer interface {
	Transform(src string, stages []t.Stage) (string, error)
	TransformCtx(ctx context.Context, src string, stages []t.Stage) (string, error)
}

func NewTransformer() Transformer {
	return &transformer{}
}

type transformer struct {
	trace_parse bool
	// trace_transform bool
}

func (i *transformer) Transform(src string, stages []t.Stage) (string, error) {
	ctx := context.Background()
	return i.TransformCtx(ctx, src, stages)
}

func (i *transformer) TransformCtx(ctx context.Context, src string, stages []t.Stage) (string, error) {
	// TEMP: unhook ctx from parser
	// _temp_parser_ctx := context.Background()
	program, err := parser.ParseCtx(ctx, src)
	if err != nil {
		return "", err
	}

	printer := node_printer.NewPrinter("TRANSFORMER_PIPELINE")
	logger := log.New(printer, "# TRANSFORMER_PIPELINE ", 0)
	ctx = logging.WithLogger(ctx, logger)

	transformer := t.NewTransformer()
	transformed_program, err := transformer.TransformCtx(ctx, program, stages)
	if transformed_program, ok := transformed_program.(*ast.Program); !ok {
		return "", fmt.Errorf("expected *t.Program, got %T", transformed_program)
	}

	if err != nil {
		return "", fmt.Errorf("transformer error: %v", err)
	}

	printer.Logf("Transformed program:")
	printer.PrintNode(transformed_program)

	return printer.String(), nil
}

var _ Transformer = &transformer{}

func Transform(src string, trace_transform bool) (string, error) {
	transformer := NewTransformer()
	return transformer.Transform(src, t.ALL_STAGES)
}

func TransformStages(src string, input interface{}, stages []t.Stage, trace_transform bool) (string, error) {
	transformer := NewTransformer()

	ctx := context.Background()
	var trace_tracer trace.Tracer
	if trace_transform {
		trace_tracer = trace.NewTracer()
		ctx = trace.WithTracer(ctx, trace_tracer)
	}

	out, err := transformer.TransformCtx(ctx, src, stages)

	if trace_tracer != nil {
		trace_tracer.Done()
		walkable, err := trace_tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		var out strings.Builder
		_ = walkable.Walk(printer.NewTracePrinter(&out, true))
		// i.stdout.Write([]byte(out.String()))
		os.Stderr.Write([]byte(out.String()))
	}

	return out, err
}
