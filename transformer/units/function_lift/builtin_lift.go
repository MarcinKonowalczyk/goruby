package function_lift

import (
	"context"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/walk"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
)

type LiftBuiltins struct {
	Statements []ast.Statement // to be lifted
}

func (f *LiftBuiltins) PostTransform(ctx context.Context, node ast.Node) ast.Node {
	switch node := node.(type) {
	case nil:
		//
	case *ast.Comment:
		//
	case *ast.ContextCallExpression:
		logging.Logf(ctx, "found context call expression \"%s\" with context %v and block %v\n",
			node.Function,
			node.Context,
			node.Block,
		)
		var replacement replacement_spec
		replacement, ok := context_call_expr_replacements[context_spec(node.Function)]
		if ok {
			logging.Logf(ctx, "replacing context call expression %s with %s\n", node.Function, replacement.name)
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
		// logging.Logf(ctx, "walking %T\n", node)
	}
	return node
}

func (f *LiftBuiltins) TransformerMarker() {}

var (
	_ walk.PostTransformerCtx = &LiftBuiltins{}
	_ walk.Transformer        = &LiftBuiltins{}
)
