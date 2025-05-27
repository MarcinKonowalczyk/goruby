package ast

import (
	"container/list"
	"fmt"
)

// A Visitor's Visit method is invoked for each node encountered by Walk. If
// the result visitor w is not nil, Walk visits each of the children of node
// with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	VisitorMarker() // marker
	Visit(node Node)
}

type Transformer interface {
	TransformerMarker() // marker
	// PreTransform(node Node) Node
	// PostTransform(node Node) Node
}

// Transforms a node *before* it is walked.
type PreTransformer interface {
	PreTransform(node Node) Node
}

// Transforms a node *after* it is walked.
type PostTransformer interface {
	PostTransform(node Node) Node
}

// Parent returns the parent node of child. If child is not found within root,
// or child does not have a parent, i.e. equals root, the bool will be false
func Parent(root, child Node) (Node, bool) {
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
	parent, ok := parentElement.Value.(Node)
	if !ok {
		return nil, false
	}
	return parent, true
}

// Path returns the path from the root of the AST till the child as a doubly
// linked list. If child is not found within root, the bool will be false and
// the list nil.
func Path(root, child Node) (*list.List, bool) {
	if !Contains(root, child) {
		return nil, false
	}
	childTree := list.New()
	l := treeToLinkedList(root)
	for e := l.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(Node)
		if !ok {
			continue
		}
		if Contains(n, child) {
			childTree.PushBack(n)
		}
	}
	return childTree, true
}

func treeToLinkedList(node Node) *list.List {
	list := list.New()
	for n := range WalkEmit(node) {
		if n != nil {
			list.PushBack(n)
		}
	}
	return list
}

type inspector func(Node)

func (f inspector) VisitorMarker() {}
func (f inspector) Visit(node Node) {
	f(node)
}

var _ Visitor = inspector(nil)

func Inspect(root Node, inspector_func func(Node)) {
	// Walk the tree and apply the filter
	Walk(root, nil, inspector(inspector_func))
}

// Contains reports whether root contains child or not. It matches child via
// pointer equality.
func Contains(root Node, child Node) bool {
	var contains bool
	filter := func(n Node) {
		if n == child {
			contains = true
		}
	}
	Walk(root, nil, inspector(filter))
	return contains
}

type visitorfunc func(Node)

func (f visitorfunc) VisitorMarker() {}
func (f visitorfunc) Visit(n Node) {
	f(n)
}

var _ Visitor = visitorfunc(nil)

// WalkEmit traverses node in depth-first order and emits each visited node
// into the channel
func WalkEmit(root Node) <-chan Node {
	out := make(chan Node)
	visitor := visitorfunc(func(n Node) {
		out <- n
	})
	go func() {
		defer close(out)
		Walk(root, nil, visitor)
	}()
	return out
}

// type Transformer interface {
// 	TransformPre(node Node) (Node, Transformer)
// 	TransformPost(node Node) (Node, Transformer)
// }

// type TransformerFunc func(Node) (Node, Transformer)

// func (f TransformerFunc) Transform(node Node) (Node, Transformer) {
// 	return f(node)
// }

// type transformer struct {
// 	pre  TransformerFunc
// 	post TransformerFunc
// }

// func (f transformer) Transform(node Node) (Node, Transformer) {
// 	if f(node) {
// 		return node, f
// 	}
// 	return nil, nil
// }

// var _ Transformer = transformer(nil)

// func WalkTransform(t Transformer, node Node) {
// 	// ...
// }

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
func Walk(
	node Node,
	transformer Transformer,
	v Visitor,
) Node {
	// if node == nil {
	// 	// we're done
	// 	return nil, post
	// }

	var pre PreTransformer
	var post PostTransformer

	if transformer != nil {
		if pre_transformer, ok := transformer.(PreTransformer); ok {
			pre = pre_transformer
		}
		if post_transformer, ok := transformer.(PostTransformer); ok {
			post = post_transformer
		}
	}

	mutating := (pre != nil || post != nil)

	// pre-transform
	if pre != nil {
		node = pre.PreTransform(node)
	}

	// visit the node
	if v != nil {
		v.Visit(node)
	}

	// walk children
	var new_node Node
	switch n := node.(type) {
	// Expressions
	case *Identifier,
		*IntegerLiteral,
		*StringLiteral,
		*SymbolLiteral,
		*Comment:
		// nothing to do

	case *FunctionLiteral:
		if mutating {
			new_params := make([]*FunctionParameter, len(n.Parameters))
			for i, x := range n.Parameters {
				new_node = Walk(x, transformer, v)
				if new_param, ok := new_node.(*FunctionParameter); ok {
					new_params[i] = new_param
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a function parameter to %T", new_param))
				}
			}
			n.Parameters = new_params
			new_node = Walk(n.Body, transformer, v)
			if new_body, ok := new_node.(*BlockStatement); ok {
				n.Body = new_body
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a function body to %T", new_body))
			}
		} else {
			for _, x := range n.Parameters {
				_ = Walk(x, transformer, v)
			}
			_ = Walk(n.Body, transformer, v)
		}

	case *FunctionParameter:
		if mutating {
			new_node = Walk(n.Default, transformer, v)
			if new_default, ok := new_node.(Expression); ok {
				n.Default = new_default
			} else {
				if new_default == nil {
					n.Default = nil
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a function parameter default to %T. Previously it was %T", new_default, n.Default))
				}
			}
		} else {
			_ = Walk(n.Default, transformer, v)
		}

	case *IndexExpression:
		if mutating {
			new_node = Walk(n.Left, transformer, v)
			if new_left, ok := new_node.(Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an index expression left to %T", new_left))
			}
			new_node = Walk(n.Index, transformer, v)
			if new_index, ok := new_node.(Expression); ok {
				n.Index = new_index
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an index expression index to %T", new_index))
			}
		} else {
			_ = Walk(n.Left, transformer, v)
			_ = Walk(n.Index, transformer, v)
		}

	case *ContextCallExpression:
		if mutating {
			new_node = Walk(n.Context, transformer, v)
			if new_node == nil {
				n.Context = nil
			} else {
				if new_context, ok := new_node.(Expression); ok {
					n.Context = new_context
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression context from %T to %T", n.Context, new_node))
				}

			}
			new_arguments := make([]Expression, len(n.Arguments))
			for i, x := range n.Arguments {
				new_node = Walk(x, transformer, v)
				if new_argument, ok := new_node.(Expression); ok {
					new_arguments[i] = new_argument
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression argument to %T", new_argument))
				}
			}
			n.Arguments = new_arguments
			if n.Block != nil {
				new_node = Walk(n.Block, transformer, v)
				if new_block, ok := new_node.(Expression); ok {
					n.Block = new_block
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a context call expression block from %T to %T", n.Block, new_node))
				}
			}
		} else {
			_ = Walk(n.Context, transformer, v)
			for _, x := range n.Arguments {
				_ = Walk(x, transformer, v)
			}
			if n.Block != nil {
				_ = Walk(n.Block, transformer, v)
			}
		}

	case *PrefixExpression:
		if mutating {
			new_node = Walk(n.Right, transformer, v)
			if new_right, ok := new_node.(Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a prefix expression right to %T", new_right))
			}

		} else {
			_ = Walk(n.Right, transformer, v)
		}

	case *InfixExpression:
		if mutating {
			new_node = Walk(n.Left, transformer, v)
			if new_left, ok := new_node.(Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an infix expression left to %T", new_left))
			}
			new_node = Walk(n.Right, transformer, v)
			if new_right, ok := new_node.(Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an infix expression right to %T", new_right))
			}
		} else {
			_ = Walk(n.Left, transformer, v)
			_ = Walk(n.Right, transformer, v)
		}

	case *MultiAssignment:
		if mutating {
			new_variables := make([]*Identifier, len(n.Variables))
			for i, x := range n.Variables {
				new_node = Walk(x, transformer, v)
				if new_variable, ok := new_node.(*Identifier); ok {
					new_variables[i] = new_variable
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a multi-assignment variable to %T", new_variable))
				}
			}
			n.Variables = new_variables
			new_values := make([]Expression, len(n.Values))
			for i, x := range n.Values {
				new_node = Walk(x, transformer, v)
				if new_value, ok := new_node.(Expression); ok {
					new_values[i] = new_value
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a multi-assignment value to %T", new_value))
				}
			}
			n.Values = new_values
		} else {
			for _, x := range n.Variables {
				_ = Walk(x, transformer, v)
			}
			for _, x := range n.Values {
				_ = Walk(x, transformer, v)
			}
		}

	case ExpressionList:
		if mutating {
			new_elements := make([]Expression, len(n))
			for i, x := range n {
				new_node = Walk(x, transformer, v)
				if new_element, ok := new_node.(Expression); ok {
					new_elements[i] = new_element
				} else {
					panic(fmt.Sprintf("ast.Walk mutated an expression list element to %T", new_element))
				}
			}
			n = new_elements
		} else {
			for _, x := range n {
				_ = Walk(x, transformer, v)
			}
		}

	// Types
	case *ArrayLiteral:
		if mutating {
			new_elements := make([]Expression, len(n.Elements))
			for i, x := range n.Elements {
				new_node = Walk(x, transformer, v)
				if new_element, ok := new_node.(Expression); ok {
					new_elements[i] = new_element
				} else {
					panic(fmt.Sprintf("ast.Walk mutated an array literal element to %T", new_element))
				}
			}
			n.Elements = new_elements
		} else {
			for _, x := range n.Elements {
				_ = Walk(x, transformer, v)
			}
		}

	case *HashLiteral:
		if mutating {
			new_map := make(map[Expression]Expression, len(n.Map))
			for k, val := range n.Map {
				new_node = Walk(k, transformer, v)
				if new_key, ok := new_node.(Expression); ok {
					new_node = Walk(val, transformer, v)
					if new_value, ok := new_node.(Expression); ok {
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
				_ = Walk(k, transformer, v)
				_ = Walk(val, transformer, v)
			}
		}

	case *ExpressionStatement:
		if mutating {
			new_node = Walk(n.Expression, transformer, v)
			// if new_expression, ok := new_node.(Expression); ok {
			// 	if new_expression != n.Expression {
			// 		fmt.Printf("# WALK *ExpressionStatement::Expression: ast.Walk mutated %T(%v) to %T(%v)\n", n.Expression, n.Expression, new_node, new_node)
			// 	}
			// 	n.Expression = new_expression
			// } else {
			// 	panic(fmt.Sprintf("ast.Walk mutated an expression statement from %T to %T", n.Expression, new_node))
			// }
			switch new_expression := new_node.(type) {
			case Expression:
				if new_expression != n.Expression {
					fmt.Printf("# WALK *ExpressionStatement::Expression: ast.Walk mutated %T(%v) to %T(%v)\n", n.Expression, n.Expression, new_node, new_node)
				}
				n.Expression = new_expression
			default:
				panic(fmt.Sprintf("ast.Walk mutated an expression statement from %T to %T", n.Expression, new_node))
			}
		} else {
			_ = Walk(n.Expression, transformer, v)
		}

	case *Assignment:
		if mutating {
			new_node = Walk(n.Left, transformer, v)
			if new_left, ok := new_node.(Expression); ok {
				if new_left != n.Left {
					fmt.Printf("*Assignment::Left: ast.Walk mutated %T to %T\n", n.Left, new_node)
				}
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an assignment left to %T", new_left))
			}
			new_node = Walk(n.Right, transformer, v)
			if new_right, ok := new_node.(Expression); ok {
				if new_right != n.Right {
					fmt.Printf("*Assignment::Right: ast.Walk mutated %T to %T\n", n.Right, new_node)
				}
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated an assignment right from %T to %T", n.Right, new_right))
			}
		} else {
			_ = Walk(n.Left, transformer, v)
			_ = Walk(n.Right, transformer, v)
		}

	case *ReturnStatement:
		if mutating {
			new_node = Walk(n.ReturnValue, transformer, v)
			if new_return_value, ok := new_node.(Expression); ok {
				if new_return_value != n.ReturnValue {
					fmt.Printf("*ReturnStatement::ReturnValue: ast.Walk mutated %T to %T\n", n.ReturnValue, new_node)
				}
				n.ReturnValue = new_return_value
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a return statement return value from %T to %T", n.ReturnValue, new_return_value))
			}

		} else {
			_ = Walk(n.ReturnValue, transformer, v)
		}

	case *BreakStatement:
		if mutating {
			new_node = Walk(n.Condition, transformer, v)
			if new_node == nil {
				n.Condition = nil
			} else {
				if new_condition, ok := new_node.(Expression); ok {
					n.Condition = new_condition
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a break statement condition from %T to %T", n.Condition, new_condition))
				}
			}
		} else {
			_ = Walk(n.Condition, transformer, v)
		}

	case *BlockStatement:
		if mutating {
			new_statements := make([]Statement, len(n.Statements))
			for i, x := range n.Statements {
				new_node = Walk(x, transformer, v)
				if new_statement, ok := new_node.(Statement); ok {
					new_statements[i] = new_statement
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a block statement to %T", new_statement))
				}
			}
			n.Statements = new_statements
		} else {
			for _, x := range n.Statements {
				_ = Walk(x, transformer, v)
			}
		}

	case *ConditionalExpression:
		if mutating {
			new_node = Walk(n.Condition, transformer, v)
			if new_condition, ok := new_node.(Expression); ok {
				n.Condition = new_condition
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a conditional expression condition to %T", new_condition))
			}
			new_node = Walk(n.Consequence, transformer, v)
			if new_consequence, ok := new_node.(*BlockStatement); ok {
				n.Consequence = new_consequence
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a conditional expression consequence to %T", new_consequence))
			}
			if n.Alternative != nil {
				new_node = Walk(n.Alternative, transformer, v)
				if new_alternative, ok := new_node.(*BlockStatement); ok {
					n.Alternative = new_alternative
				} else {
					panic(fmt.Sprintf("ast.Walk mutated a conditional expression alternative to %T", new_alternative))
				}
			}
		} else {
			_ = Walk(n.Condition, transformer, v)
			_ = Walk(n.Consequence, transformer, v)
			if n.Alternative != nil {
				_ = Walk(n.Alternative, transformer, v)
			}
		}

	case *LoopExpression:
		if mutating {
			// new_node = Walk(n.Condition, transformer, v)
			// if new_condition, ok := new_node.(Expression); ok {
			// n.Condition = new_condition
			// } else {
			// panic(fmt.Sprintf("ast.Walk mutated a loop expression condition to %T", new_condition))
			// }
			new_node = Walk(n.Block, transformer, v)
			if new_block, ok := new_node.(*BlockStatement); ok {
				n.Block = new_block
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a loop expression block from %T to %T", n.Block, new_block))
			}
		} else {
			// _ = Walk(n.Condition, transformer, v)
			_ = Walk(n.Block, transformer, v)
		}

	// Program
	case *Program:
		if mutating {
			new_statements := make([]Statement, len(n.Statements))
			for i, statement := range n.Statements {
				new_node = Walk(statement, transformer, v)
				if new_node == nil {
					new_statements[i] = nil
				} else {
					if new_statement, ok := new_node.(Statement); ok {
						new_statements[i] = new_statement
					} else {
						panic(fmt.Sprintf("ast.Walk mutated a program statement from %T to %T", statement, new_statement))
					}
				}
			}
			n.Statements = new_statements
		} else {
			for _, x := range n.Statements {
				_ = Walk(x, transformer, v)
			}
		}

	case *RangeLiteral:
		if mutating {
			new_node = Walk(n.Left, transformer, v)
			if new_left, ok := new_node.(Expression); ok {
				n.Left = new_left
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a range literal left from %T to %T", n.Left, new_left))
			}
			new_node = Walk(n.Right, transformer, v)
			if new_right, ok := new_node.(Expression); ok {
				n.Right = new_right
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a range literal right from %T to %T", n.Right, new_right))
			}
		} else {
			_ = Walk(n.Left, transformer, v)
			_ = Walk(n.Right, transformer, v)
		}

	case nil:
		// nothing to do

	case *Splat:
		if mutating {
			new_node = Walk(n.Value, transformer, v)
			if new_value, ok := new_node.(Expression); ok {
				n.Value = new_value
			} else {
				panic(fmt.Sprintf("ast.Walk mutated a splat value from %T to %T", n.Value, new_value))
			}
		} else {
			_ = Walk(n.Value, transformer, v)
		}

	case *FloatLiteral:
		// nothing to do

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	// post-transform
	// NOTE: we will *not* visit the result of the post-transform
	//       that's kinda the point
	if post != nil {
		// fmt.Println("pre-transform", node)
		node = post.PostTransform(node)
		// fmt.Println("post-transform", node)
	}
	// v.Visit(nil)

	return node
}
