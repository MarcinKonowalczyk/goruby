package builtin_lifts

import (
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

type Lift struct {
	Statements []ast.Statement // to be lifted
}

func (f *Lift) PostTransform(node ast.Node) ast.Node {
	switch node := node.(type) {
	case nil:
		//
	case *ast.Comment:
		//
	case *ast.ContextCallExpression:
		// fmt.Printf("%v (%T)\n", node, node)
		fmt.Printf("# TRANSFORMER found context call expression \"%s\"\n", node.Function)
		// f.functions = append(f.functions, node)

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
		fmt.Printf("# TRANSFORMER walking %T\n", node)
	}
	return node
}

func (f *Lift) TransformerMarker() {}

var (
	_ ast.PostTransformer = &Lift{}
	_ ast.Transformer     = &Lift{}
)
