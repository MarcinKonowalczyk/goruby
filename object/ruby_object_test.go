package object

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/utils"
)

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

		_, err := function.Call(context, nil)
		utils.AssertNoError(t, err)

		var expected ast.Node = functionBody
		utils.Assert(t, reflect.DeepEqual(expected, actualEvalNode), "Expected Eval argument to equal\n%v\n\tgot\n%v\n", expected, actualEvalNode)
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

		_, err := function.Call(context, nil)
		utils.Assert(t, reflect.DeepEqual(evalErr, err), "Expected error to equal\n%v\n\tgot\n%v\n", evalErr, err)
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

		_, err := function.Call(context, nil)
		utils.AssertNoError(t, err)

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

			_, err := function.Call(context, nil, NewInteger(300), NewString("sym"))
			utils.AssertNoError(t, err)

			{
				expected := NewInteger(300)
				actual, ok := evalEnv.Get("foo")
				utils.Assert(t, ok, "Expected function parameter %q to be in Eval env", "foo")
				utils.AssertEqualCmpAny(t, expected, actual, CompareRubyObjectsForTests)
			}
			{
				expected := NewString("sym")
				actual, ok := evalEnv.Get("bar")
				utils.Assert(t, ok, "Expected function parameter %q to be in Eval env", "bar")
				utils.AssertEqualCmpAny(t, expected, actual, CompareRubyObjectsForTests)
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

			_, err := function.Call(context, nil, NewInteger(300), NewSymbol("sym"))
			utils.AssertNoError(t, err)

			{
				actual, ok := evalEnv.Get("foo")
				utils.Assert(t, ok, "Expected function parameter %q to be in Eval env", "foo")
				utils.AssertEqualCmpAny(t, NewInteger(12), actual, CompareRubyObjectsForTests)
			}
			{
				actual, ok := evalEnv.Get("bar")
				utils.Assert(t, ok, "Expected function parameter %q to be in Eval env", "bar")
				utils.AssertEqualCmpAny(t, NewInteger(300), actual, CompareRubyObjectsForTests)
			}
			{
				actual, ok := evalEnv.Get("qux")
				utils.Assert(t, ok, "Expected function parameter %q to be in Eval env", "qux")
				utils.AssertEqualCmpAny(t, NewSymbol("sym"), actual, CompareRubyObjectsForTests)
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

			result, _ := function.Call(context, nil)
			utils.AssertEqualCmpAny(t, NewInteger(8), result, CompareRubyObjectsForTests)
		})
		t.Run("wrapped into a return value", func(t *testing.T) {
			context := &callContext{
				env:  NewMainEnvironment(),
				eval: func(ast.Node, Environment) (RubyObject, error) { return &ReturnValue{Value: NewInteger(8)}, nil },
			}

			function := &Function{}

			result, _ := function.Call(context, nil)
			utils.AssertEqualCmpAny(t, NewInteger(8), result, CompareRubyObjectsForTests)
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
			_, err := function.Call(context, nil, NewString("foo"))
			utils.AssertError(t, err, NewWrongNumberOfArgumentsError(0, 1))
		})

		t.Run("with default arguments", func(t *testing.T) {
			function.Parameters = []*FunctionParameter{
				{Name: "x", Default: TRUE},
				{Name: "y"},
			}

			_, err := function.Call(context, nil, NewInteger(8))
			utils.AssertNoError(t, err)
		})
	})
}
