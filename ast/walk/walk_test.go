package walk_test

import (
	"container/list"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/walk"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
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
