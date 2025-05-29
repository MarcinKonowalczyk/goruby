package transformer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
)

func Transform(node *ast.Program, stages []Stage, trace_transform bool) (*ast.Program, error) {
	// new context with stdout logger
	logger := log.New(os.Stdout, "# TRANSFORMER ", 0)
	ctx := logging.WithLogger(context.Background(), logger)

	t := &transformer{}
	if trace_transform {
		t.tracer = trace.NewTracer()
	}
	transformed_node, err := t.TransformCtx(ctx, node, stages)
	if err != nil {
		return nil, err
	}
	node, ok := transformed_node.(*ast.Program)
	if !ok {
		return nil, fmt.Errorf("expected *ast.Program, got %T", transformed_node)
	}
	return node, nil
}

func NewTransformer() *transformer {
	return &transformer{
		tracer: trace.NewTracer(),
	}
}

type Transformer interface {
	Transform(node ast.Node, stages []Stage) (ast.Node, error)
	TransformCtx(ctx context.Context, node ast.Node, stages []Stage) (ast.Node, error)
}

var (
	_ Transformer = &transformer{}
)
