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

func (f *LiftBuiltins) Visit(ctx context.Context, node ast.Node) (ast.Node, walk.Flag) {
	switch node := node.(type) {
	case *ast.ContextCallExpression:
		walk.WalkChildrenCtx(ctx, node, f)
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
			}, walk.SKIP
		}
		return walk.NOOP, walk.SKIP
	default:
		// continue walking
		return walk.NOOP, walk.WALK
	}
}

func (f *LiftBuiltins) TransformerMarker() {}

var (
	_ walk.Visitor = &LiftBuiltins{}
)
