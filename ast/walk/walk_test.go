package walk_test

import (
	"container/list"
	"context"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/assert"
	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/walk"

	"github.com/brunoga/deep"
)

func Test_Parent(t *testing.T) {
	t.Run("parent found", func(t *testing.T) {
		child := &ast.Assignment{
			Left:  &ast.Identifier{Value: "x"},
			Right: &ast.IntegerLiteral{Value: 2},
		}
		parent := &ast.ExpressionStatement{Expression: child}
		root := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 3},
				},
				parent,
			},
		}

		p, ok := walk.Parent(root, child)

		assert.That(t, ok, "Expected child to be contained within root, was not")
		assert.That(t, reflect.DeepEqual(parent, p), "Expected parent to equal\n%+#v\n\tgot\n%+#v\n", parent, p)
	})
	t.Run("parent is root", func(t *testing.T) {
		root := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 3},
				},
			},
		}

		_, ok := walk.Parent(root, root)

		assert.That(t, !ok, "Expected bool to return false")
	})
	t.Run("child not found", func(t *testing.T) {
		root := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 3},
				},
				&ast.ExpressionStatement{
					Expression: &ast.Assignment{
						Left:  &ast.Identifier{Value: "x"},
						Right: &ast.IntegerLiteral{Value: 2},
					},
				},
			},
		}

		_, ok := walk.Parent(root, &ast.IntegerLiteral{Value: 3})

		assert.That(t, !ok, "Expected child not to be contained within root")
	})
}

func Test_Path(t *testing.T) {
	t.Run("child found", func(t *testing.T) {
		child := &ast.Assignment{
			Left:  &ast.Identifier{Value: "x"},
			Right: &ast.IntegerLiteral{Value: 2},
		}
		root := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 3},
				},
				&ast.ExpressionStatement{
					Expression: child,
				},
			},
		}

		path, ok := walk.Path(root, child)

		assert.That(t, ok, "Expected child to be contained within root, was not")

		expected := list.New()
		expected.PushBack(root)
		expected.PushBack(root.Statements[1])
		expected.PushBack(child)

		assert.That(t,
			reflect.DeepEqual(expected, path),
			"Expected AST path to equal\n%+#v\n\tgot\n%+#v\n", expected, path,
		)
	})
	t.Run("child not found", func(t *testing.T) {
		root := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 3},
				},
				&ast.ExpressionStatement{
					Expression: &ast.Assignment{
						Left:  &ast.Identifier{Value: "x"},
						Right: &ast.IntegerLiteral{Value: 2},
					},
				},
			},
		}

		_, ok := walk.Path(root, &ast.IntegerLiteral{Value: 3})
		assert.That(t, !ok, "Expected child not to be contained within root")
	})
}

func Test_treeToList(t *testing.T) {
	root := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.IntegerLiteral{Value: 3},
			},
			&ast.ExpressionStatement{
				Expression: &ast.Assignment{
					Left:  &ast.Identifier{Value: "x"},
					Right: &ast.IntegerLiteral{Value: 2},
				},
			},
		},
	}

	actual := walk.TreeToLinkedList(root)

	expected := list.New()
	expected.PushBack(root)
	expected.PushBack(&ast.ExpressionStatement{
		Expression: &ast.IntegerLiteral{Value: 3},
	})
	expected.PushBack(&ast.IntegerLiteral{Value: 3})
	expected.PushBack(&ast.ExpressionStatement{
		Expression: &ast.Assignment{
			Left:  &ast.Identifier{Value: "x"},
			Right: &ast.IntegerLiteral{Value: 2},
		},
	})
	expected.PushBack(&ast.Assignment{
		Left:  &ast.Identifier{Value: "x"},
		Right: &ast.IntegerLiteral{Value: 2},
	})
	expected.PushBack(&ast.Identifier{Value: "x"})
	expected.PushBack(&ast.IntegerLiteral{Value: 2})

	assert.That(t,
		reflect.DeepEqual(expected, actual),
		"Expected list to equal\n%+#v\n\tgot\n%+#v\n", expected, actual,
	)
}

func Test_Contains(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		needle := &ast.IntegerLiteral{Value: 1}
		statement := &ast.ExpressionStatement{
			Expression: needle,
		}

		assert.That(t,
			walk.Contains(statement, needle),
			"Expected node to be within statement, was not",
		)
	})
	t.Run("not the same pointer", func(t *testing.T) {
		needle := &ast.IntegerLiteral{Value: 1}
		other_needle := &ast.IntegerLiteral{Value: 1}
		statement := &ast.ExpressionStatement{
			Expression: needle,
		}

		assert.That(t,
			!walk.Contains(statement, other_needle),
			"Expected node not to be within statement",
		)

	})
	t.Run("not found", func(t *testing.T) {
		needle := &ast.IntegerLiteral{Value: 3}
		statement := &ast.ExpressionStatement{
			Expression: &ast.Assignment{
				Left:  &ast.Identifier{Value: "foo"},
				Right: &ast.StringLiteral{Value: "bar"},
			},
		}

		assert.That(t,
			!walk.Contains(statement, needle),
			"Expected node not to be within statement",
		)
	})
}

// Make a tree
// ```
// 3
// x = 2
// "hello"
// ```
func makeTestTree() (*ast.Program, []ast.Node) {
	tree := &ast.Program{
		Statements: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: &ast.IntegerLiteral{Value: 3},
			},
			&ast.ExpressionStatement{
				Expression: &ast.Assignment{
					Left:  &ast.Identifier{Value: "x"},
					Right: &ast.IntegerLiteral{Value: 2},
				},
			},
			&ast.ExpressionStatement{
				Expression: &ast.StringLiteral{Value: "hello"},
			},
		},
	}

	tree2 := deep.MustCopy(tree)

	// nodes in walk order.
	expected := []ast.Node{
		tree2,
		tree2.Statements[0],
		tree2.Statements[0].(*ast.ExpressionStatement).Expression,
		tree2.Statements[1],
		tree2.Statements[1].(*ast.ExpressionStatement).Expression.(*ast.Assignment),
		tree2.Statements[1].(*ast.ExpressionStatement).Expression.(*ast.Assignment).Left,
		tree2.Statements[1].(*ast.ExpressionStatement).Expression.(*ast.Assignment).Right,
		tree2.Statements[2],
		tree2.Statements[2].(*ast.ExpressionStatement).Expression.(*ast.StringLiteral),
	}

	return tree, expected
}
func compareTrees(a, b []ast.Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		at := reflect.TypeOf(a[i])
		bt := reflect.TypeOf(b[i])
		if at != bt {
			return false
		}

		if at == reflect.TypeOf(&ast.IntegerLiteral{}) {
			// Compare IntegerLiteral values
			if a[i].(*ast.IntegerLiteral).Value != b[i].(*ast.IntegerLiteral).Value {
				return false
			}
		}
		if at == reflect.TypeOf(&ast.StringLiteral{}) {
			// Compare StringLiteral values
			if a[i].(*ast.StringLiteral).Value != b[i].(*ast.StringLiteral).Value {
				return false
			}
		}
		// ...

	}
	return true
}
func Test_WalkCtx(t *testing.T) {
	t.Run("walks through all nodes", func(t *testing.T) {
		root, expected := makeTestTree()

		var visited []ast.Node
		walk.WalkCtx(nil, root, walk.VisitorFunc(func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			visited = append(visited, n)
			return walk.NOOP, walk.WALK
		}))

		assert.EqualCmp(t, expected, visited, compareTrees)
	})

	t.Run("walks does not continue", func(t *testing.T) {
		root, expected := makeTestTree()

		var visited []ast.Node
		walk.WalkCtx(nil, root, walk.VisitorFunc(func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			visited = append(visited, n)
			if _, ok := n.(*ast.Assignment); ok {
				return walk.NOOP, walk.SKIP
			}
			return walk.NOOP, walk.WALK
		}))

		// pop the children of the assignment from the expected list
		var assignment_index int = 0
		for i, n := range expected {
			if _, ok := n.(*ast.Assignment); ok {
				assignment_index = i
				break
			}
		}
		assert.That(t, assignment_index > 0, "Expected assignment to be found in expected list")
		new_expected := []ast.Node{}
		// remove two nodes after the assignment, but not the assignment itself
		new_expected = append(new_expected, expected[:assignment_index+1]...)
		new_expected = append(new_expected, expected[assignment_index+3:]...)

		assert.EqualCmp(t, new_expected, visited, compareTrees)
	})
}

func Test_WalkCtxTransform(t *testing.T) {
	t.Run("transforms int literal value", func(t *testing.T) {
		root, expected := makeTestTree()
		var visited []ast.Node
		walk.WalkCtx(nil, root, walk.VisitorFunc(func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			visited = append(visited, n)
			if np, ok := n.(*ast.IntegerLiteral); ok {
				np.Value = 42 // Transform IntegerLiteral value to 42
			}
			return walk.NOOP, walk.WALK
		}))

		// Set expected IntegerLiteral values to 42
		for _, n := range expected {
			if il, ok := n.(*ast.IntegerLiteral); ok {
				il.Value = 42 // Set expected IntegerLiteral values to 42
			}
		}

		assert.EqualCmp(t, expected, visited, compareTrees)
	})
	t.Run("transforms node type", func(t *testing.T) {
		root, expected := makeTestTree()
		var visited []ast.Node

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			if _, ok := (n).(*ast.IntegerLiteral); ok {
				n2 := &ast.StringLiteral{Value: "42"}
				visited = append(visited, n2)
				return n2, walk.SKIP
			} else {
				visited = append(visited, n)
			}
			return walk.NOOP, walk.WALK
		}

		walk.WalkCtx(nil, root, walk.VisitorFunc(vf))

		// Set expected
		for i, n := range expected {
			if _, ok := n.(*ast.IntegerLiteral); ok {
				expected[i] = &ast.StringLiteral{Value: "42"} // Set expected IntegerLiteral to StringLiteral
			}
		}

		assert.EqualCmp(t, expected, visited, compareTrees)
	})

	t.Run("deletes node but fails", func(t *testing.T) {
		root, _ := makeTestTree()

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			if _, ok := n.(*ast.IntegerLiteral); ok {
				return walk.DELETE, walk.SKIP
			}
			return walk.NOOP, walk.WALK
		}

		assert.Equal(t, root.Code(), "3; x = 2; \"hello\"")

		assert.Panic(t, func() {
			walk.WalkCtx(nil, root, walk.VisitorFunc(vf))
		}, func(t testing.TB, rec any) {
			assert.ContainsString(t, assert.Type[string](t, rec), "cannot delete right side of assignment")
		})

		assert.Equal(t, root.Code(), "3; x = 2; \"hello\"")
	})

	t.Run("deletes node", func(t *testing.T) {
		root, _ := makeTestTree()

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			if _, ok := n.(*ast.Assignment); ok {
				return walk.DELETE, walk.SKIP
			}
			return walk.NOOP, walk.WALK
		}

		assert.Equal(t, root.Code(), "3; x = 2; \"hello\"")
		walk.WalkCtx(nil, root, walk.VisitorFunc(vf))
		assert.Equal(t, root.Code(), "3; \"hello\"")
	})

	t.Run("delete walks into its own marker", func(t *testing.T) {
		root, _ := makeTestTree()

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			if _, ok := n.(*ast.Assignment); ok {
				return walk.DELETE, walk.WALK
			}
			return walk.NOOP, walk.WALK
		}

		assert.Panic(t, func() { walk.WalkCtx(nil, root, walk.VisitorFunc(vf)) },
			func(t testing.TB, rec any) {
				assert.ContainsString(t, assert.Type[string](t, rec), "walked into a delete node")
			})
	})

	t.Run("walk inner-transform", func(t *testing.T) {
		root, _ := makeTestTree()

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			if _, ok := n.(*ast.Assignment); ok {
				// walk the assignment node
				vf2 := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
					if il, ok := n.(*ast.IntegerLiteral); ok {
						il.Value = 42 // Transform IntegerLiteral value to 42
					}
					return n, walk.WALK
				}
				return walk.WalkCtx(nil, n, walk.VisitorFunc(vf2)), walk.WALK
			}
			return n, walk.WALK
		}

		assert.Equal(t, root.Code(), "3; x = 2; \"hello\"")
		walk.WalkCtx(nil, root, walk.VisitorFunc(vf))
		assert.Equal(t, root.Code(), "3; x = 42; \"hello\"")
	})

	t.Run("walk inner-transform dont continue", func(t *testing.T) {
		root, expected := makeTestTree()

		outer_visited := []ast.Node{}
		inner_visited := []ast.Node{}

		vf := func(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
			outer_visited = append(outer_visited, n)
			if _, ok := n.(*ast.Assignment); ok {
				// vf2 := func(n ast.Node) (ast.Node, bool) {
				// 	inner_visited = append(inner_visited, n)
				// 	return nil, true
				// }
				// walk.WalkCtx(nil, n, walk.VisitorFunc(vf2))
				walk.Inspect(n, func(n ast.Node) {
					inner_visited = append(inner_visited, n)
				})
				return walk.NOOP, walk.SKIP
			}
			return n, walk.WALK
		}

		walk.WalkCtx(nil, root, walk.VisitorFunc(vf))

		// pop the children of the assignment from the expected list
		var assignment_index int = 0
		for i, n := range expected {
			if _, ok := n.(*ast.Assignment); ok {
				assignment_index = i
				break
			}
		}
		assert.That(t, assignment_index > 0, "Expected assignment to be found in expected list")
		new_expected := []ast.Node{}
		// remove two nodes after the assignment, but not the assignment itself
		new_expected = append(new_expected, expected[:assignment_index+1]...)
		new_expected = append(new_expected, expected[assignment_index+3:]...)

		assert.EqualCmp(t, new_expected, outer_visited, compareTrees)
		assert.EqualCmp(t, expected[assignment_index:assignment_index+3], inner_visited, compareTrees)
	})

	t.Run("walk recursive", func(t *testing.T) {
		root, expected := makeTestTree()

		walker := &recursiveWalker{}
		walk.WalkCtx(nil, root, walker)

		assert.That(t, len(walker.visited) > 0, "Expected walker to visit at least one node")
		assert.EqualCmp(t, expected, walker.visited, compareTrees)
	})

}

type recursiveWalker struct {
	visited []ast.Node
}

func (w *recursiveWalker) Visit(_ context.Context, n ast.Node) (ast.Node, walk.Flag) {
	w.visited = append(w.visited, n)
	switch n := n.(type) {
	case *ast.Program:
		walk.WalkChildrenCtx(nil, n, w)
		return walk.NOOP, walk.SKIP
	default:
		return walk.NOOP, walk.WALK
	}
}
