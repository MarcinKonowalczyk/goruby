package function_lift

import (
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/walk"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
)

type call struct {
	name       string
	parameters []string
}

func (c *call) String() string {
	return fmt.Sprintf("call(name=%s, parameters=%v)", c.name, c.parameters)
}

type LiftBlocks struct {
	call_stack []*call
	// functions which we lifted and want to now stick in the top level of the program
	LiftedFunctions []*ast.FunctionLiteral
	// change the names of some of the functions
	NameChanges map[string]string // map of old name to new name
	// Pass number. this is a two-Pass transformation.
	Pass int
}

func (f *LiftBlocks) Visit(ctx context.Context, node ast.Node) (ast.Node, walk.Flag) {
	switch node := node.(type) {
	case *ast.FunctionLiteral:
		// increment the call stack
		parameter_names := make([]string, len(node.Parameters))
		for i, param := range node.Parameters {
			parameter_names[i] = param.Name
		}
		f.call_stack = append(f.call_stack, &call{name: node.Name, parameters: parameter_names})

		// f.Visit(node) // walk
		walk.WalkChildrenCtx(ctx, node, f)

		// decrement the call stack
		if len(f.call_stack) > 0 {
			f.call_stack = f.call_stack[:len(f.call_stack)-1]
		}

		// call the post-transform function
		if f.Pass == 0 {
			return f.transformFunctionLiteralPass0(ctx, node), walk.SKIP
		}

		return walk.NOOP, walk.SKIP

	case *ast.ContextCallExpression:
		if f.Pass == 1 {
			walk.WalkChildrenCtx(ctx, node, f)
			return f.transformContextCallExpressionPass1(ctx, node), walk.SKIP
		}

	}
	return nil, true
}

func (f *LiftBlocks) transformFunctionLiteralPass0(ctx context.Context, node *ast.FunctionLiteral) ast.Node {
	logging.Logf(ctx, "lifting function '%s' with parameters '%v'\n", node.Name, node.Parameters)
	// logging.Logf(ctx, "call stack: %v\n", f.call_stack)
	if len(f.call_stack) > 0 {
		parent := f.call_stack[len(f.call_stack)-1]

		// And add all the parameters from the parent
		extra_parameters := make([]*ast.FunctionParameter, 0)

		// TODO:
		// NOTE: This won't really work that well. What if we define a local variable
		// and need to capture that?

		for _, param := range parent.parameters {
			function_parameter := &ast.FunctionParameter{Name: param}
			extra_parameters = append(extra_parameters, function_parameter)
		}

		var lifted_function *ast.FunctionLiteral
		var replacement_node *ast.IndexExpression
		if len(extra_parameters) > 0 {
			// we need a factory function
			factory_body := &ast.Statements{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: node,
					},
				},
			}
			lifted_function = &ast.FunctionLiteral{
				Name:       node.Name,
				Parameters: extra_parameters,
				Body:       factory_body,
			}
			var index ast.ExpressionList
			for _, param := range parent.parameters {
				index = append(index, &ast.Identifier{Value: param})
			}
			replacement_node = &ast.IndexExpression{
				Left:  &ast.Identifier{Value: node.Name},
				Index: index,
			}
		} else {
			// lift the function as is just the node
			lifted_function = node
			var index ast.ExpressionList
			for _, param := range node.Parameters {
				index = append(index, &ast.Identifier{Value: param.Name})
			}
			for _, param := range parent.parameters {
				index = append(index, &ast.Identifier{Value: param})
			}
			replacement_node = &ast.IndexExpression{
				Left:  &ast.Identifier{Value: node.Name},
				Index: index,
			}
		}

		f.LiftedFunctions = append(f.LiftedFunctions, lifted_function)

		return replacement_node
	} else {
		old_name := node.Name
		new_name := fmt.Sprintf("__lifted_%s", node.Name)
		logging.Logf(ctx, "lifting top-level function '%s' to '%s'\n", old_name, new_name)

		node.Name = new_name
		if f.NameChanges == nil {
			f.NameChanges = make(map[string]string)
		}
		f.NameChanges[old_name] = new_name
		f.LiftedFunctions = append(f.LiftedFunctions, node)

		return walk.DELETE
	}
}

func (f *LiftBlocks) transformContextCallExpressionPass1(ctx context.Context, node *ast.ContextCallExpression) ast.Node {
	if node.Context == nil {
		new_name, ok := f.NameChanges[node.Function]
		if ok {
			new_node := &ast.IndexExpression{
				Left:  &ast.Identifier{Value: new_name},
				Index: ast.ExpressionList(node.Arguments),
			}
			logging.Logf(ctx, "Replacing context call expression '%s' with '%s'\n", node.Function, new_name)
			return new_node
		}
	}
	return walk.NOOP
}

func (f *LiftBlocks) TransformerMarker() {}

var (
	_ walk.Visitor = &LiftBlocks{}
)
