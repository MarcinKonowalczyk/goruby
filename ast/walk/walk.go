package walk

import (
	"container/list"
	"context"
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/trace"
)

// A Visitor's Visit method is invoked for each node encountered by Walk. If
// the result visitor w is not nil, Walk visits each of the children of node
// with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(ctx context.Context, node ast.Node) (ast.Node, Flag)
}

// type Transformer interface {
// 	TransformerMarker() // marker
// 	// PreTransform(node ast.Node) ast.Node
// 	// PreTransform(ctx context.Context, node ast.Node) ast.Node
// 	// PostTransform(node ast.Node) ast.Node
// 	// PostTransform(ctx context.Context, node ast.Node) ast.Node
// }

// // Transforms a node *before* it is walked.
// type PreTransformer interface {
// 	PreTransform(node ast.Node) ast.Node
// }

// type PreTransformerCtx interface {
// 	PreTransform(ctx context.Context, node ast.Node) ast.Node
// }

// // Transforms a node *after* it is walked.
// type PostTransformer interface {
// 	PostTransform(node ast.Node) ast.Node
// }

// type PostTransformerCtx interface {
// 	PostTransform(ctx context.Context, node ast.Node) ast.Node
// }

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
func (f inspector) Visit(ctx context.Context, node ast.Node) (ast.Node, Flag) {
	f(node)
	return NOOP, WALK // continue walking, don't replace the node
}

var _ Visitor = inspector(nil)

func Inspect(root ast.Node, inspector_func func(ast.Node)) {
	// Walk the tree and apply the filter
	WalkCtx(nil, root, inspector(inspector_func))
}

// Contains reports whether root contains child or not. It matches child via
// pointer equality.
func Contains(root ast.Node, child ast.Node) bool {
	var contains bool
	filter := func(ctx context.Context, n ast.Node) (ast.Node, Flag) {
		if n == child {
			contains = true
		}
		return nil, Flag(!contains) // continue walking until we find the child
	}
	WalkCtx(nil, root, VisitorFunc(filter))
	return contains
}

type VisitorFunc func(context.Context, ast.Node) (ast.Node, Flag)

func (f VisitorFunc) VisitorMarker() {}
func (f VisitorFunc) Visit(ctx context.Context, n ast.Node) (ast.Node, Flag) {
	return f(ctx, n)
}

var _ Visitor = VisitorFunc(nil)

// WalkEmit traverses node in depth-first order and emits each visited node
// into the channel
func WalkEmit(root ast.Node) <-chan ast.Node {
	out := make(chan ast.Node)
	visitor := VisitorFunc(func(ctx context.Context, n ast.Node) (ast.Node, Flag) {
		out <- n
		return nil, true // continue walking
	})
	go func() {
		defer close(out)
		WalkCtx(nil, root, visitor)
	}()
	return out
}

type Flag bool

const (
	WALK Flag = true
	SKIP Flag = false
)

var (
	DELETE ast.Node = &ast.SpecialNode{Value: "delete"}
	NOOP   ast.Node = nil
)

// var WALK_CHILDREN bool = true
// var SKIP_CHILDREN bool = false

// Pre transformer with context which wraps a PreTransformer without context
// type pre_transformer_ctx struct {
// 	transformer PreTransformer
// }

// func (p pre_transformer_ctx) PreTransform(_ context.Context, node ast.Node) ast.Node {
// 	if p.transformer == nil {
// 		return node
// 	}
// 	return p.transformer.PreTransform(node)
// }

// var _ PreTransformerCtx = pre_transformer_ctx{}

// Post transformer with context which wraps a PostTransformer without context
// type post_transformer_ctx struct {
// 	transformer PostTransformer
// }

// func (p post_transformer_ctx) PostTransform(_ context.Context, node ast.Node) ast.Node {
// 	if p.transformer == nil {
// 		return node
// 	}
// 	return p.transformer.PostTransform(node)
// }

// var _ PostTransformerCtx = post_transformer_ctx{}

type nodeReplacement struct {
	index int
	node  ast.Node
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
func WalkCtx(
	ctx context.Context,
	node ast.Node,
	v Visitor,
) ast.Node {
	defer trace.TraceCtx(ctx, fmt.Sprintf("%T", node))()
	flag := WALK

	// visit the node
	var new_node ast.Node = nil
	new_node, flag = v.Visit(ctx, node)
	var any_new_children bool = false

	if new_node != nil {
		any_new_children = true
		node = new_node
	}

	if flag == SKIP {
		trace.MessageCtx(ctx, "skipping children")
		if any_new_children {
			return node // return the node with new children
		} else {
			// no new children, return nil to indicate no changes
			return NOOP
		}
	}

	return WalkChildrenCtx(ctx, node, v)
}

// walks children of node, applying the visitor v to each child.
func WalkChildrenCtx(
	ctx context.Context,
	node ast.Node,
	v Visitor,
) ast.Node {
	var any_children_changed bool = false
	// walk children
	switch n := node.(type) {
	case *ast.Identifier,
		*ast.FloatLiteral,
		*ast.IntegerLiteral,
		*ast.StringLiteral,
		*ast.SymbolLiteral,
		*ast.Comment,
		nil:
		// nothing to do, these are leaf nodes

	case *ast.FunctionLiteral:
		for _, x := range n.Parameters {
			WalkCtx(ctx, x, v)
		}
		WalkCtx(ctx, n.Body, v)

	case *ast.FunctionParameter:
		WalkCtx(ctx, n.Default, v)

	case *ast.IndexExpression:
		WalkCtx(ctx, n.Left, v)
		WalkCtx(ctx, n.Index, v)

	case *ast.ContextCallExpression:

		// walk context
		nn := WalkCtx(ctx, n.Context, v)
		switch nn := nn.(type) {
		case nil:
			// no changes, continue
		case ast.Expression:
			// if we have a new expression, we replace the old one
			any_children_changed = true
			n.Context = nn
		case *ast.SpecialNode:
			if nn.Value == "delete" {
				// for now we disallow deleting context, but we might want to allow it in the future..??
				panic("cannot delete context of context call expression")
			} else {
				panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
			}
		default:
			panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
		}

		// walk arguments
		replacements := make([]nodeReplacement, 0)
		deletions := make([]int, 0)
		for i, x := range n.Arguments {
			nn := WalkCtx(ctx, x, v)
			switch nn := nn.(type) {
			case nil:
				// no changes, continue
			case *ast.SpecialNode:
				if nn.Value == "delete" {
					trace.MessageCtx(ctx, "deleting argument")
					deletions = append(deletions, i)
				} else {
					panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
				}
			case ast.Expression:
				// if we have a new expression, we replace the old one
				replacements = append(replacements, nodeReplacement{
					index: i,
					node:  nn,
				})
			default:
				panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
			}
		}

		// NOTE: we do replacements first since they don't clobber the indices
		if len(replacements) > 0 {
			any_children_changed = true
			for _, r := range replacements {
				n.Arguments[r.index] = r.node.(ast.Expression)
			}
		}
		if len(deletions) > 0 {
			any_children_changed = true
			n.Arguments = deleteIndices(n.Arguments, deletions)
		}

		// walk block
		nn = WalkCtx(ctx, n.Block, v)
		switch nn := nn.(type) {
		case nil:
			// no changes, continue
		case *ast.SpecialNode:
			if nn.Value == "delete" {
				trace.MessageCtx(ctx, "deleting block")
				any_children_changed = true
				n.Block = nil // delete the block
			} else {
				panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
			}
		case ast.Expression:
			any_children_changed = true
			n.Block = nn
		default:
			panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
		}

	case *ast.PrefixExpression:
		WalkCtx(ctx, n.Right, v)

	case *ast.InfixExpression:
		WalkCtx(ctx, n.Left, v)
		WalkCtx(ctx, n.Right, v)

	case *ast.MultiAssignment:
		for _, x := range n.Variables {
			WalkCtx(ctx, x, v)
		}
		for _, x := range n.Values {
			WalkCtx(ctx, x, v)
		}

	case ast.ExpressionList:
		for _, x := range n {
			WalkCtx(ctx, x, v)
		}

	case *ast.ArrayLiteral:
		for _, x := range n.Elements {
			WalkCtx(ctx, x, v)
		}

	case *ast.HashLiteral:
		for k, val := range n.Map {
			WalkCtx(ctx, k, v)
			WalkCtx(ctx, val, v)
		}

	case *ast.Assignment:
		nn := WalkCtx(ctx, n.Left, v)
		switch nn := nn.(type) {
		case nil:
			// no changes
		case ast.Expression:
			// if we have a new expression, we replace the old one
			any_children_changed = true
			n.Left = nn
		case *ast.SpecialNode:
			if nn.Value == "delete" {
				panic("cannot delete left side of assignment")
			} else {
				panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
			}
		default:
			panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
		}
		nn = WalkCtx(ctx, n.Right, v)
		switch nn := nn.(type) {
		case nil:
			// no changes, continue
		case ast.Expression:
			// if we have a new expression, we replace the old one
			any_children_changed = true
			n.Right = nn
		case *ast.SpecialNode:
			if nn.Value == "delete" {
				panic("cannot delete right side of assignment")
			} else {
				panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
			}
		default:
			panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
		}

	case *ast.ReturnStatement:
		WalkCtx(ctx, n.ReturnValue, v)

	case *ast.BreakStatement:
		WalkCtx(ctx, n.Condition, v)

	case *ast.Statements:
		for _, x := range n.Statements {
			WalkCtx(ctx, x, v)
		}

	case *ast.ConditionalExpression:
		WalkCtx(ctx, n.Condition, v)
		WalkCtx(ctx, n.Consequence, v)
		if n.Alternative != nil {
			WalkCtx(ctx, n.Alternative, v)
		}

	case *ast.LoopExpression:
		WalkCtx(ctx, n.Block, v)

	// Program
	case *ast.Program:
		var replacements []nodeReplacement
		var deletions []int

		for i, x := range n.Statements {
			nn := WalkCtx(ctx, x, v)
			switch nn := nn.(type) {
			case nil:
				// no changes, continue
			case *ast.SpecialNode:
				if nn.Value == "delete" {
					trace.MessageCtx(ctx, "deleting")
					deletions = append(deletions, i)
				} else {
					panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
				}
			case ast.Statement:
				// if we have a new statement, we replace the old one
				any_children_changed = true
				replacements = append(replacements, nodeReplacement{
					index: i,
					node:  nn,
				})
			default:
				panic(fmt.Sprintf("ast.Walk: expected ast.Statement, got %T", nn))
			}
		}

		// NOTE: we do replacements fist since they don't clobber the indices
		if len(replacements) > 0 {
			for _, r := range replacements {
				n.Statements[r.index] = r.node.(ast.Statement)
			}
		}
		if len(deletions) > 0 {
			n.Statements = deleteIndices(n.Statements, deletions)
		}

	case *ast.ExpressionStatement:
		nn := WalkCtx(ctx, n.Expression, v)
		switch nn := nn.(type) {
		case nil:
			// no changes, continue
		case *ast.SpecialNode:
			if nn.Value == "delete" {
				trace.MessageCtx(ctx, "deleting")
				return DELETE // return the special delete node
			} else {
				panic(fmt.Sprintf("ast.Walk: unexpected special node %s", nn.Value))
			}
		case ast.Expression:
			any_children_changed = true
			n.Expression = nn
		default:
			panic(fmt.Sprintf("ast.Walk: expected ast.Expression, got %T", nn))
		}

	case *ast.RangeLiteral:
		WalkCtx(ctx, n.Left, v)
		WalkCtx(ctx, n.Right, v)

	case *ast.Splat:
		WalkCtx(ctx, n.Value, v)

	case *ast.SpecialNode:
		if n.Value == "delete" {
			panic("walked into a delete node. did you return `walk.DeleteNode, true` from your visitor? (should be `false` instead)")
		} else {
			panic(fmt.Sprintf("ast.Walk: unexpected special node %s", n.Value))
		}
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	if any_children_changed {
		return node // return the node with new children
	} else {
		// no new children, return nil to indicate no changes
		return NOOP
	}
}

// func Walk(
// 	node ast.Node,
// 	v Visitor,
// ) {
// 	WalkCtx(context.Background(), node, v)
// }
