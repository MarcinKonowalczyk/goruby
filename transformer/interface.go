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

func Transform(node *ast.Program, trace_transform bool) (*ast.Program, error) {
	// new context with stdout logger
	logger := log.New(os.Stdout, "# TRANSFORMER ", 0)
	ctx := logging.WithLogger(context.Background(), logger)

	t := NewTransformer().(*transformer)
	if trace_transform {
		t.tracer = trace.NewTracer()
	}
	transformed_node, err := t.TransformCtx(ctx, node)
	if err != nil {
		return nil, err
	}
	node, ok := transformed_node.(*ast.Program)
	if !ok {
		return nil, fmt.Errorf("expected *ast.Program, got %T", transformed_node)
	}
	return node, nil
}

type Transformer interface {
	Transform(node ast.Node) (ast.Node, error)
}

func NewTransformer() Transformer {
	return &transformer{}
}

var (
	_ Transformer = &transformer{}
)
