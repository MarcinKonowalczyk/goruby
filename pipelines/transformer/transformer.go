package transformer

import (
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/printer"
	t "github.com/MarcinKonowalczyk/goruby/transformer"
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
	program, err := parser.Parse(src)
	// if tracer != nil {
	// 	walkable, err := tracer.ToWalkable()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	walkable.Walk(trace_printer.NewTracePrinter())
	// }
	if err != nil {
		return "", err
	}

	transformer := t.NewTransformer()
	transformed_program, err := transformer.TransformCtx(ctx, program, stages)
	if transformed_program, ok := transformed_program.(*ast.Program); !ok {
		return "", fmt.Errorf("expected *t.Program, got %T", transformed_program)
	}

	if err != nil {
		return "", fmt.Errorf("transformer error: %v", err)
	}

	printer := printer.NewPrinter("TRANSFORM_PIPELINE")
	printer.PrintNode(transformed_program)

	return printer.String(), nil
}

var _ Transformer = &transformer{}

func Transform(src string) (string, error) {
	transformer := NewTransformer()
	return transformer.Transform(src, t.ALL_STAGES)
}

// func TransformWithComments(filename string, input interface{}) (string, error) {
// }

func TransformStages(src string, input interface{}, stages []t.Stage) (string, error) {
	transformer := NewTransformer()
	ctx := context.Background()
	// logger := log.New(os.Stdout, "# TRANSFORMER ", 0)
	// ctx = logging.WithLogger(ctx, logger)
	return transformer.TransformCtx(ctx, src, t.ALL_STAGES)
	// transformer := NewTransformer()
	// return transformer.Transform(filename, input, stages)
}
