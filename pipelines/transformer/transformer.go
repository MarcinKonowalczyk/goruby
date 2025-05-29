package transformer

import (
	"fmt"
	"go/token"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/printer"
	trace_printer "github.com/MarcinKonowalczyk/goruby/trace/printer"
	t "github.com/MarcinKonowalczyk/goruby/transformer"
)

type Transformer interface {
	Transform(filename string, input interface{}, stages []t.Stage) (string, error)
}

func NewTransformer() Transformer {
	return &transformer{}
}

type transformer struct {
	trace_parse bool
	// trace_transform bool
}

func (i *transformer) Transform(filename string, input interface{}, stages []t.Stage) (string, error) {
	program, tracer, err := parser.ParseFileEx(token.NewFileSet(), filename, input, i.trace_parse)
	if tracer != nil {
		walkable, err := tracer.ToWalkable()
		if err != nil {
			panic(err)
		}
		walkable.Walk(trace_printer.NewTracePrinter())
	}
	if err != nil {
		return "", err
	}

	transformer := t.NewTransformer()
	transformed_program, err := transformer.Transform(program, stages)
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

func Transform(filename string, input interface{}) (string, error) {
	transformer := NewTransformer()
	return transformer.Transform(filename, input, t.ALL_STAGES)
}

func TransformStages(filename string, input interface{}, stages []t.Stage) (string, error) {
	transformer := NewTransformer()
	return transformer.Transform(filename, input, stages)
}
