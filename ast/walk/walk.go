package walk

import (
	"container/list"
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

// A Visitor's Visit method is invoked for each node encountered by Walk. If
// the result visitor w is not nil, Walk visits each of the children of node
// with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	VisitorMarker() // marker
	Visit(node ast.Node)
}

type Transformer interface {
	TransformerMarker() // marker
	// PreTransform(node ast.Node) ast.Node
	// PreTransform(ctx context.Context, node ast.Node) ast.Node
	// PostTransform(node ast.Node) ast.Node
	// PostTransform(ctx context.Context, node ast.Node) ast.Node
}

// Transforms a node *before* it is walked.
type PreTransformer interface {
	PreTransform(node ast.Node) ast.Node
}

type PreTransformerCtx interface {
	PreTransform(ctx context.Context, node ast.Node) ast.Node
}

// Transforms a node *after* it is walked.
type PostTransformer interface {
	PostTransform(node ast.Node) ast.Node
}

type PostTransformerCtx interface {
	PostTransform(ctx context.Context, node ast.Node) ast.Node
}

// Parent returns the parent node of child. If child is not found within root,
// or child does not have a parent, i.e. equals root, the bool will be false
func Parent(root, child ast.Node) (ast.Node, bool) {
	if root == child {
		return nil, false
	}
	if !Contains(root, child) {
		return nil, false
	}
	path, ok := Path(root, child)
	if !ok {
		return nil, false
	}
	parentElement := path.Back().Prev()
	if parentElement == nil {
		return nil, false
	}
	parent, ok := parentElement.Value.(ast.Node)
	if !ok {
		return nil, false
	}
	return parent, true
}

// Path returns the path from the root of the AST till the child as a doubly
// linked list. If child is not found within root, the bool will be false and
// the list nil.
func Path(root, child ast.Node) (*list.List, bool) {
	if !Contains(root, child) {
		return nil, false
	}
	childTree := list.New()
	l := TreeToLinkedList(root)
	for e := l.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(ast.Node)
		if !ok {
			continue
		}
		if Contains(n, child) {
			childTree.PushBack(n)
		}
	}
	return childTree, true
}

func TreeToLinkedList(node ast.Node) *list.List {
	list := list.New()
	for n := range WalkEmit(node) {
		if n != nil {
			list.PushBack(n)
		}
	}
	return list
}

type inspector func(ast.Node)

func (f inspector) VisitorMarker() {}
func (f inspector) Visit(node ast.Node) {
	f(node)
}

var _ Visitor = inspector(nil)

func Inspect(root ast.Node, inspector_func func(ast.Node)) {
	// Walk the tree and apply the filter
	Walk(root, nil, inspector(inspector_func))
}

// Contains reports whether root contains child or not. It matches child via
// pointer equality.
func Contains(root ast.Node, child ast.Node) bool {
	var contains bool
	filter := func(n ast.Node) {
		if n == child {
			contains = true
		}
	}
	Walk(root, nil, inspector(filter))
	return contains
}

type visitorfunc func(ast.Node)

func (f visitorfunc) VisitorMarker() {}
func (f visitorfunc) Visit(n ast.Node) {
	f(n)
}

var _ Visitor = visitorfunc(nil)

// WalkEmit traverses node in depth-first order and emits each visited node
// into the channel
func WalkEmit(root ast.Node) <-chan ast.Node {
	out := make(chan ast.Node)
	visitor := visitorfunc(func(n ast.Node) {
		out <- n
	})
	go func() {
		defer close(out)
		Walk(root, nil, visitor)
	}()
	return out
}

// type Transformer interface {
// 	TransformPre(node ast.Node) (ast.Node, Transformer)
// 	TransformPost(node ast.Node) (ast.Node, Transformer)
// }

// type TransformerFunc func(ast.Node) (ast.Node, Transformer)

// func (f TransformerFunc) Transform(node ast.Node) (ast.Node, Transformer) {
// 	return f(node)
// }

// type transformer struct {
// 	pre  TransformerFunc
// 	post TransformerFunc
// }

// func (f transformer) Transform(node ast.Node) (ast.Node, Transformer) {
// 	if f(node) {
// 		return node, f
// 	}
// 	return nil, nil
// }

// var _ Transformer = transformer(nil)

// func WalkTransform(t Transformer, node ast.Node) {
// 	// ...
// }

// Pre transformer with context which wraps a PreTransformer without context
type pre_transformer_ctx struct {
	transformer PreTransformer
}

func (p pre_transformer_ctx) PreTransform(_ context.Context, node ast.Node) ast.Node {
	if p.transformer == nil {
		return node
	}
	return p.transformer.PreTransform(node)
}

var _ PreTransformerCtx = pre_transformer_ctx{}

// Post transformer with context which wraps a PostTransformer without context
type post_transformer_ctx struct {
	transformer PostTransformer
}

func (p post_transformer_ctx) PostTransform(_ context.Context, node ast.Node) ast.Node {
	if p.transformer == nil {
		return node
	}
	return p.transformer.PostTransform(node)
}

var _ PostTransformerCtx = post_transformer_ctx{}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
func WalkCtx(
	ctx context.Context,
	node ast.Node,
	transformer Transformer,
	v Visitor,
) ast.Node {
	// if node == nil {
	// 	// we're done
	// 	return nil, post
	// }

	var pre PreTransformerCtx
	var post PostTransformerCtx

	if transformer != nil {
		switch pre_transformer := transformer.(type) {
		case PreTransformer:
			pre = pre_transformer_ctx{transformer: pre_transformer}
		case PreTransformerCtx:
			pre = pre_transformer
		default:
			// do nothing
		}
		switch post_transformer := transformer.(type) {
		case PostTransformer:
			post = post_transformer_ctx{transformer: post_transformer}
		case PostTransformerCtx:
			post = post_transformer
		default:
			// do nothing
		}
	}

	mutating := (pre != nil || post != nil)

	// pre-transform
	if pre != nil {
		node = pre.PreTransform(ctx, node)
	}

	// visit the node
	if v != nil {
		v.Visit(node)
	}

	// walk children
	var new_node ast.Node
	switch n := node.(type) {
	// Expressions
	case *ast.Identifier,
		*ast.IntegerLiteral,
		*ast.StringLiteral,
		*ast.SymbolLiteral,
		*ast.Comment:
		// nothing to do

	case *ast.FunctionLiteral:
		if mutating {
			new_params := make([]*ast.FunctionParameter, len(n.Parameters))
			for i, x := range n.Parameters {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_param, ok := new_node.(*ast.FunctionParameter); ok {
					new_params[i] = new_param
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a function parameter to %T", new_param))
				}
			}
			n.Parameters = new_params
			new_node = WalkCtx(ctx, n.Body, transformer, v)
			if new_body, ok := new_node.(*ast.BlockStatement); ok {
				n.Body = new_body
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a function body to %T", new_body))
			}
		} else {
			for _, x := range n.Parameters {
				_ = WalkCtx(ctx, x, transformer, v)
			}
			_ = WalkCtx(ctx, n.Body, transformer, v)
		}

	case *ast.FunctionParameter:
		if mutating {
			new_node = WalkCtx(ctx, n.Default, transformer, v)
			if new_default, ok := new_node.(ast.Expression); ok {
				n.Default = new_default
			} else {
				if new_default == nil {
					n.Default = nil
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a function parameter default to %T. Previously it was %T", new_default, n.Default))
				}
			}
		} else {
			_ = WalkCtx(ctx, n.Default, transformer, v)
		}

	case *ast.IndexExpression:
		if mutating {
			new_node = WalkCtx(ctx, n.Left, transformer, v)
			if new_left, ok := new_node.(ast.Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an index expression left to %T", new_left))
			}
			new_node = WalkCtx(ctx, n.Index, transformer, v)
			if new_index, ok := new_node.(ast.Expression); ok {
				n.Index = new_index
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an index expression index to %T", new_index))
			}
		} else {
			_ = WalkCtx(ctx, n.Left, transformer, v)
			_ = WalkCtx(ctx, n.Index, transformer, v)
		}

	case *ast.ContextCallExpression:
		if mutating {
			new_node = WalkCtx(ctx, n.Context, transformer, v)
			if new_node == nil {
				n.Context = nil
			} else {
				if new_context, ok := new_node.(ast.Expression); ok {
					n.Context = new_context
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression context from %T to %T", n.Context, new_node))
				}

			}
			new_arguments := make([]ast.Expression, len(n.Arguments))
			for i, x := range n.Arguments {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_argument, ok := new_node.(ast.Expression); ok {
					new_arguments[i] = new_argument
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression argument to %T", new_argument))
				}
			}
			n.Arguments = new_arguments
			if n.Block != nil {
				new_node = WalkCtx(ctx, n.Block, transformer, v)
				if new_block, ok := new_node.(ast.Expression); ok {
					n.Block = new_block
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression block from %T to %T", n.Block, new_node))
				}
			}
		} else {
			_ = WalkCtx(ctx, n.Context, transformer, v)
			for _, x := range n.Arguments {
				_ = WalkCtx(ctx, x, transformer, v)
			}
			if n.Block != nil {
				_ = WalkCtx(ctx, n.Block, transformer, v)
			}
		}

	case *ast.PrefixExpression:
		if mutating {
			new_node = WalkCtx(ctx, n.Right, transformer, v)
			if new_right, ok := new_node.(ast.Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a prefix expression right to %T", new_right))
			}

		} else {
			_ = WalkCtx(ctx, n.Right, transformer, v)
		}

	case *ast.InfixExpression:
		if mutating {
			new_node = WalkCtx(ctx, n.Left, transformer, v)
			if new_left, ok := new_node.(ast.Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an infix expression left to %T", new_left))
			}
			new_node = WalkCtx(ctx, n.Right, transformer, v)
			if new_right, ok := new_node.(ast.Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an infix expression right to %T", new_right))
			}
		} else {
			_ = WalkCtx(ctx, n.Left, transformer, v)
			_ = WalkCtx(ctx, n.Right, transformer, v)
		}

	case *ast.MultiAssignment:
		if mutating {
			new_variables := make([]*ast.Identifier, len(n.Variables))
			for i, x := range n.Variables {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_variable, ok := new_node.(*ast.Identifier); ok {
					new_variables[i] = new_variable
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a multi-assignment variable to %T", new_variable))
				}
			}
			n.Variables = new_variables
			new_values := make([]ast.Expression, len(n.Values))
			for i, x := range n.Values {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_value, ok := new_node.(ast.Expression); ok {
					new_values[i] = new_value
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a multi-assignment value to %T", new_value))
				}
			}
			n.Values = new_values
		} else {
			for _, x := range n.Variables {
				_ = WalkCtx(ctx, x, transformer, v)
			}
			for _, x := range n.Values {
				_ = WalkCtx(ctx, x, transformer, v)
			}
		}

	case ast.ExpressionList:
		if mutating {
			new_elements := make([]ast.Expression, len(n))
			for i, x := range n {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_element, ok := new_node.(ast.Expression); ok {
					new_elements[i] = new_element
				} else {
					panic(fmt.Sprintf("ast.Walk mutated an expression list element to %T", new_element))
				}
			}
			n = new_elements
		} else {
			for _, x := range n {
				_ = WalkCtx(ctx, x, transformer, v)
			}
		}

	// Types
	case *ast.ArrayLiteral:
		if mutating {
			new_elements := make([]ast.Expression, len(n.Elements))
			for i, x := range n.Elements {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_element, ok := new_node.(ast.Expression); ok {
					new_elements[i] = new_element
				} else {
					panic(fmt.Sprintf("ast.Walk mutated an array literal element to %T", new_element))
				}
			}
			n.Elements = new_elements
		} else {
			for _, x := range n.Elements {
				_ = WalkCtx(ctx, x, transformer, v)
			}
		}

	case *ast.HashLiteral:
		if mutating {
			new_map := make(map[ast.Expression]ast.Expression, len(n.Map))
			for k, val := range n.Map {
				new_node = WalkCtx(ctx, k, transformer, v)
				if new_key, ok := new_node.(ast.Expression); ok {
					new_node = WalkCtx(ctx, val, transformer, v)
					if new_value, ok := new_node.(ast.Expression); ok {
						new_map[new_key] = new_value
					} else {
						panic(fmt.Sprintf("ast.Walk mutated a hash literal value to %T", new_value))
					}
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a hash literal key to %T", new_key))
				}
			}
			n.Map = new_map
		} else {
			for k, val := range n.Map {
				_ = WalkCtx(ctx, k, transformer, v)
				_ = WalkCtx(ctx, val, transformer, v)
			}
		}

	case *ast.ExpressionStatement:
		if mutating {
			new_node = WalkCtx(ctx, n.Expression, transformer, v)
			// if new_expression, ok := new_node.(ast.Expression); ok {
			// 	if new_expression != n.Expression {
			// 		fmt.Printf("# WALK *ExpressionStatement::Expression: ast.Walk mutated %T(%v) to %T(%v)\n", n.Expression, n.Expression, new_node, new_node)
			// 	}
			// 	n.Expression = new_expression
			// } else {
			// 	panic(fmt.Sprintf("ast.Walk mutated an expression statement from %T to %T", n.Expression, new_node))
			// }
			switch new_expression := new_node.(type) {
			case ast.Expression:
				if new_expression != n.Expression {
					fmt.Printf("# WALK *ExpressionStatement::Expression: ast.Walk mutated %T(%v) to %T(%v)\n", n.Expression, n.Expression, new_node, new_node)
				}
				n.Expression = new_expression
			case nil:
				// we got nil. remove the expression statement
				n = nil
			default:
				panic(fmt.Sprintf("ast.Walk mutated an expression statement from %T to %T", n.Expression, new_node))
			}
		} else {
			_ = WalkCtx(ctx, n.Expression, transformer, v)
		}

	case *ast.Assignment:
		if mutating {
			new_node = WalkCtx(ctx, n.Left, transformer, v)
			if new_left, ok := new_node.(ast.Expression); ok {
				if new_left != n.Left {
					fmt.Printf("*Assignment::Left: ast.Walk mutated %T to %T\n", n.Left, new_node)
				}
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an assignment left to %T", new_left))
			}
			new_node = WalkCtx(ctx, n.Right, transformer, v)
			if new_right, ok := new_node.(ast.Expression); ok {
				if new_right != n.Right {
					fmt.Printf("*Assignment::Right: ast.Walk mutated %T to %T\n", n.Right, new_node)
				}
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an assignment right from %T to %T", n.Right, new_right))
			}
		} else {
			_ = WalkCtx(ctx, n.Left, transformer, v)
			_ = WalkCtx(ctx, n.Right, transformer, v)
		}

	case *ast.ReturnStatement:
		if mutating {
			new_node = WalkCtx(ctx, n.ReturnValue, transformer, v)
			if new_return_value, ok := new_node.(ast.Expression); ok {
				if new_return_value != n.ReturnValue {
					fmt.Printf("*ReturnStatement::ReturnValue: ast.Walk mutated %T to %T\n", n.ReturnValue, new_node)
				}
				n.ReturnValue = new_return_value
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a return statement return value from %T to %T", n.ReturnValue, new_return_value))
			}

		} else {
			_ = WalkCtx(ctx, n.ReturnValue, transformer, v)
		}

	case *ast.BreakStatement:
		if mutating {
			new_node = WalkCtx(ctx, n.Condition, transformer, v)
			if new_node == nil {
				n.Condition = nil
			} else {
				if new_condition, ok := new_node.(ast.Expression); ok {
					n.Condition = new_condition
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a break statement condition from %T to %T", n.Condition, new_condition))
				}
			}
		} else {
			_ = WalkCtx(ctx, n.Condition, transformer, v)
		}

	case *ast.BlockStatement:
		if mutating {
			new_statements := make([]ast.Statement, len(n.Statements))
			for i, x := range n.Statements {
				new_node = WalkCtx(ctx, x, transformer, v)
				if new_statement, ok := new_node.(ast.Statement); ok {
					new_statements[i] = new_statement
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a block statement to %T", new_statement))
				}
			}
			n.Statements = new_statements
		} else {
			for _, x := range n.Statements {
				_ = WalkCtx(ctx, x, transformer, v)
			}
		}

	case *ast.ConditionalExpression:
		if mutating {
			new_node = WalkCtx(ctx, n.Condition, transformer, v)
			if new_condition, ok := new_node.(ast.Expression); ok {
				n.Condition = new_condition
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a conditional expression condition to %T", new_condition))
			}
			new_node = WalkCtx(ctx, n.Consequence, transformer, v)
			if new_consequence, ok := new_node.(*ast.BlockStatement); ok {
				n.Consequence = new_consequence
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a conditional expression consequence to %T", new_consequence))
			}
			if n.Alternative != nil {
				new_node = WalkCtx(ctx, n.Alternative, transformer, v)
				if new_alternative, ok := new_node.(*ast.BlockStatement); ok {
					n.Alternative = new_alternative
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a conditional expression alternative to %T", new_alternative))
				}
			}
		} else {
			_ = WalkCtx(ctx, n.Condition, transformer, v)
			_ = WalkCtx(ctx, n.Consequence, transformer, v)
			if n.Alternative != nil {
				_ = WalkCtx(ctx, n.Alternative, transformer, v)
			}
		}

	case *ast.LoopExpression:
		if mutating {
			// new_node = WalkCtx(ctx, n.Condition, transformer, v)
			// if new_condition, ok := new_node.(ast.Expression); ok {
			// n.Condition = new_condition
			// } else {
			// panic(fmt.Sprintf("ast.Walk mutated a loop expression condition to %T", new_condition))
			// }
			new_node = WalkCtx(ctx, n.Block, transformer, v)
			if new_block, ok := new_node.(*ast.BlockStatement); ok {
				n.Block = new_block
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a loop expression block from %T to %T", n.Block, new_block))
			}
		} else {
			// _ = WalkCtx(ctx, n.Condition, transformer, v)
			_ = WalkCtx(ctx, n.Block, transformer, v)
		}

	// Program
	case *ast.Program:
		if mutating {
			new_statements := make([]ast.Statement, len(n.Statements))
			for i, statement := range n.Statements {
				new_node = WalkCtx(ctx, statement, transformer, v)
				if new_node == nil {
					new_statements[i] = nil
				} else {
					if new_statement, ok := new_node.(ast.Statement); ok {
						new_statements[i] = new_statement
					} else {
						panic(fmt.Sprintf("ast.Walk mutated a program statement from %T to %T", statement, new_statement))
					}
				}
			}
			n.Statements = new_statements
		} else {
			for _, x := range n.Statements {
				_ = WalkCtx(ctx, x, transformer, v)
			}
		}

	case *ast.RangeLiteral:
		if mutating {
			new_node = WalkCtx(ctx, n.Left, transformer, v)
			if new_left, ok := new_node.(ast.Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a range literal left from %T to %T", n.Left, new_left))
			}
			new_node = WalkCtx(ctx, n.Right, transformer, v)
			if new_right, ok := new_node.(ast.Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a range literal right from %T to %T", n.Right, new_right))
			}
		} else {
			_ = WalkCtx(ctx, n.Left, transformer, v)
			_ = WalkCtx(ctx, n.Right, transformer, v)
		}

	case nil:
		// nothing to do

	case *ast.Splat:
		if mutating {
			new_node = WalkCtx(ctx, n.Value, transformer, v)
			if new_value, ok := new_node.(ast.Expression); ok {
				n.Value = new_value
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a splat value from %T to %T", n.Value, new_value))
			}
		} else {
			_ = WalkCtx(ctx, n.Value, transformer, v)
		}

	case *ast.FloatLiteral:
		// nothing to do

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	// post-transform
	// NOTE: we will *not* visit the result of the post-transform
	//       that's kinda the point
	if post != nil {
		// fmt.Println("pre-transform", node)
		node = post.PostTransform(ctx, node)
		// fmt.Println("post-transform", node)
	}
	// v.Visit(nil)

	return node
}

func Walk(
	node ast.Node,
	transformer Transformer,
	v Visitor,
) ast.Node {
	return WalkCtx(context.Background(), node, transformer, v)
}
