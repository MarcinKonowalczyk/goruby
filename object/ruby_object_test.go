package object

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/utils"
)

func mustCall(t *testing.T) func(obj RubyObject, err error) RubyObject {
	return func(obj RubyObject, err error) RubyObject {
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.FailNow()
		}
		return obj
	}
}

func TestFunctionCall(t *testing.T) {
	t.Run("calls CallContext#Eval with its Body", func(t *testing.T) {
		functionBody := &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 13},
				},
			},
		}

		function := &Function{
			Body: functionBody,
		}

		var actualEvalNode ast.Node
		context := &callContext{
			env: NewMainEnvironment(),
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				actualEvalNode = node
				return nil, nil
			},
		}

		mustCall(t)(function.Call(context))

		var expected ast.Node = functionBody
		if !reflect.DeepEqual(expected, actualEvalNode) {
			t.Logf("Expected Eval argument to equal\n%v\n\tgot\n%v\n", expected, actualEvalNode)
			t.Fail()
		}
	})
	t.Run("returns any error returned by CallContext#Eval", func(t *testing.T) {
		evalErr := fmt.Errorf("An error")

		context := &callContext{
			env:  NewMainEnvironment(),
			eval: func(ast.Node, Environment) (RubyObject, error) { return nil, evalErr },
		}

		function := &Function{
			Parameters: []*FunctionParameter{},
		}

		_, err := function.Call(context)

		if !reflect.DeepEqual(evalErr, err) {
			t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", evalErr, err)
			t.Fail()
		}
	})
	t.Run("uses the function env as env for CallContext#Eval", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", NewInteger(42))
		contextEnv.Set("bar", NewString("not reachable in Eval"))
		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		functionEnv := NewEnvironment()
		functionEnv.Set("foo", NewSymbol("bar"))
		function := &Function{
			Parameters: []*FunctionParameter{},
			Env:        functionEnv,
		}

		mustCall(t)(function.Call(context))

		{
			expected := NewSymbol("bar")
			actual, ok := evalEnv.Get("foo")

			utils.Assert(t, ok, "Expected key 'foo' to be in Eval env")
			utils.AssertEqualCmpAny(t, expected, actual, CompareRubyObjectsForTests)
		}

		_, ok := evalEnv.Get("bar")
		utils.Assert(t, !ok, "Expected key 'bar' not to be in Eval env")

	})
	t.Run("puts the Call args into the env for CallContext#Eval", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", NewInteger(42))
		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		t.Run("without default params", func(t *testing.T) {
			function := &Function{
				Parameters: []*FunctionParameter{
					{Name: "foo"},
					{Name: "bar"},
				},
			}

			mustCall(t)(function.Call(context, NewInteger(300), NewString("sym")))

			{
				expected := NewInteger(300)
				actual, ok := evalEnv.Get("foo")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "foo")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := NewString("sym")
				actual, ok := evalEnv.Get("bar")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "bar")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
		})
		t.Run("with default params", func(t *testing.T) {
			t.Skip()
			function := &Function{
				Parameters: []*FunctionParameter{
					{Name: "foo", Default: NewInteger(12)},
					{Name: "bar"},
					{Name: "qux"},
				},
			}

			mustCall(t)(function.Call(context, NewInteger(300), NewSymbol("sym")))

			{
				expected := NewInteger(12)
				actual, ok := evalEnv.Get("foo")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "foo")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := NewInteger(300)
				actual, ok := evalEnv.Get("bar")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "bar")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := NewSymbol("sym")
				actual, ok := evalEnv.Get("qux")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "qux")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
		})
	})
	t.Run("returns the object returned by CallContext#Eval", func(t *testing.T) {
		t.Run("vanilla object", func(t *testing.T) {
			context := &callContext{
				env:  NewMainEnvironment(),
				eval: func(ast.Node, Environment) (RubyObject, error) { return NewInteger(8), nil },
			}

			function := &Function{}

			result, _ := function.Call(context)

			expected := NewInteger(8)

			if !reflect.DeepEqual(expected, result) {
				t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, result)
				t.Fail()
			}
		})
		t.Run("wrapped into a return value", func(t *testing.T) {
			context := &callContext{
				env:  NewMainEnvironment(),
				eval: func(ast.Node, Environment) (RubyObject, error) { return &ReturnValue{Value: NewInteger(8)}, nil },
			}

			function := &Function{}

			result, _ := function.Call(context)

			expected := NewInteger(8)

			if !reflect.DeepEqual(expected, result) {
				t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, result)
				t.Fail()
			}
		})
	})
	t.Run("validates that the arguments match the function parameters", func(t *testing.T) {
		context := &callContext{
			env:  NewMainEnvironment(),
			eval: func(ast.Node, Environment) (RubyObject, error) { return nil, nil },
		}

		function := &Function{
			Parameters: []*FunctionParameter{},
		}

		t.Run("without block argument", func(t *testing.T) {
			expected := NewWrongNumberOfArgumentsError(0, 1)

			_, err := function.Call(context, NewString("foo"))

			if !reflect.DeepEqual(expected, err) {
				t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", expected, err)
				t.Fail()
			}
		})

		t.Run("with default arguments", func(t *testing.T) {
			function.Parameters = []*FunctionParameter{
				{Name: "x", Default: TRUE},
				{Name: "y"},
			}

			_, err := function.Call(context, NewInteger(8))

			if err != nil {
				t.Logf("Expected no error, got %T:%v\n", err, err)
				t.Fail()
			}
		})
	})
}
