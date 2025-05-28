package builtin_lifts

import (
	"context"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
)

type Lift struct {
	Statements []ast.Statement // to be lifted
}

func (f *Lift) PostTransform(ctx context.Context, node ast.Node) ast.Node {
	switch node := node.(type) {
	case nil:
		//
	case *ast.Comment:
		//
	case *ast.ContextCallExpression:
		logging.Logf(ctx, "found context call expression \"%s\" with context %s and block %s\n",
			node.Function, node.Context, node.Block)
		var replacement replacement_spec
		replacement, ok := context_call_expr_replacements[context_spec(node.Function)]
		if ok {
			// found a replacement !!
			f.Statements = append(f.Statements, replacement.statement)
			return &ast.IndexExpression{
				Left: &ast.Identifier{Value: replacement.name},
				Index: ast.ExpressionList{
					node.Context,
					node.Block,
				},
			}
		}
	default:
		logging.Logf(ctx, "walking %T\n", node)
	}
	return node
}

func (f *Lift) TransformerMarker() {}

var (
	_ ast.PostTransformerCtx = &Lift{}
	_ ast.Transformer        = &Lift{}
)
