package object_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
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

		function := &object.Function{
			Body: functionBody,
		}

		var actualEvalNode ast.Node
		ctx := call.NewContext(context.Background(), nil, object.NewMainEnvironment())
		ctx = call.WithEval(ctx, func(node ast.Node, env env.Environment[ruby.Object]) (ruby.Object, error) {
			actualEvalNode = node
			return nil, nil
		})

		_, err := function.Call(ctx)
		assert.NoError(t, err)

		var expected ast.Node = functionBody
		assert.That(t, reflect.DeepEqual(expected, actualEvalNode), "Expected Eval argument to equal\n%v\n\tgot\n%v\n", expected, actualEvalNode)
	})
	t.Run("returns any error returned by CallContext#Eval", func(t *testing.T) {
		evalErr := fmt.Errorf("An error")
		ctx := call.NewContext(context.Background(), nil, object.NewMainEnvironment())
		ctx = call.WithEval(ctx, func(ast.Node, env.Environment[ruby.Object]) (ruby.Object, error) {
			return nil, evalErr
		})

		function := &object.Function{
			Parameters: []*object.FunctionParameter{},
		}

		_, err := function.Call(ctx)
		assert.That(t, reflect.DeepEqual(evalErr, err), "Expected error to equal\n%v\n\tgot\n%v\n", evalErr, err)
	})
	t.Run("uses the function env as env for CallContext#Eval", func(t *testing.T) {
		contextEnv := env.NewEnvironment[ruby.Object]()
		contextEnv.Set("bar", object.NewString("not reachable in Eval"))
		var evalEnv env.Environment[ruby.Object]
		ctx := call.NewContext(context.Background(), nil, contextEnv)
		ctx = call.WithEval(ctx, func(node ast.Node, env env.Environment[ruby.Object]) (ruby.Object, error) {
			evalEnv = env
			return nil, nil
		})

		functionEnv := env.NewEnvironment[ruby.Object]()
		functionEnv.Set("foo", object.NewSymbol("bar"))
		function := &object.Function{
			Parameters: []*object.FunctionParameter{},
			Env:        functionEnv,
		}

		_, err := function.Call(ctx)
		assert.NoError(t, err)

		{
			expected := object.NewSymbol("bar")
			actual, ok := evalEnv.Get("foo")
			assert.That(t, ok, "Expected key 'foo' to be in Eval env")
			assert.EqualCmpAny(t, expected, actual, object.CompareRubyObjectsForTests)
		}

		_, ok := evalEnv.Get("bar")
		assert.That(t, !ok, "Expected key 'bar' not to be in Eval env")

	})
	t.Run("puts the Call args into the env for CallContext#Eval", func(t *testing.T) {
		contextEnv := env.NewEnvironment[ruby.Object]()
		var evalEnv env.Environment[ruby.Object]
		ctx := call.NewContext(context.Background(), nil, contextEnv)
		ctx = call.WithEval(ctx, func(node ast.Node, env env.Environment[ruby.Object]) (ruby.Object, error) {
			evalEnv = env
			return nil, nil
		})

		t.Run("without default params", func(t *testing.T) {
			function := &object.Function{
				Parameters: []*object.FunctionParameter{
					{Name: "foo"},
					{Name: "bar"},
				},
			}

			_, err := function.Call(ctx, object.NewInteger(300), object.NewString("sym"))
			assert.NoError(t, err)

			{
				expected := object.NewInteger(300)
				actual, ok := evalEnv.Get("foo")
				assert.That(t, ok, "Expected function parameter %q to be in Eval env", "foo")
				assert.EqualCmpAny(t, expected, actual, object.CompareRubyObjectsForTests)
			}
			{
				expected := object.NewString("sym")
				actual, ok := evalEnv.Get("bar")
				assert.That(t, ok, "Expected function parameter %q to be in Eval env", "bar")
				assert.EqualCmpAny(t, expected, actual, object.CompareRubyObjectsForTests)
			}
		})
		t.Run("with default params", func(t *testing.T) {
			t.Skip()
			function := &object.Function{
				Parameters: []*object.FunctionParameter{
					{Name: "foo", Default: object.NewInteger(12)},
					{Name: "bar"},
					{Name: "qux"},
				},
			}

			_, err := function.Call(ctx, object.NewInteger(300), object.NewSymbol("sym"))
			assert.NoError(t, err)

			{
				actual, ok := evalEnv.Get("foo")
				assert.That(t, ok, "Expected function parameter %q to be in Eval env", "foo")
				assert.EqualCmpAny(t, object.NewInteger(12), actual, object.CompareRubyObjectsForTests)
			}
			{
				actual, ok := evalEnv.Get("bar")
				assert.That(t, ok, "Expected function parameter %q to be in Eval env", "bar")
				assert.EqualCmpAny(t, object.NewInteger(300), actual, object.CompareRubyObjectsForTests)
			}
			{
				actual, ok := evalEnv.Get("qux")
				assert.That(t, ok, "Expected function parameter %q to be in Eval env", "qux")
				assert.EqualCmpAny(t, object.NewSymbol("sym"), actual, object.CompareRubyObjectsForTests)
			}
		})
	})
	t.Run("returns the object returned by CallContext#Eval", func(t *testing.T) {
		t.Run("vanilla object", func(t *testing.T) {
			ctx := call.NewContext(context.Background(), nil, object.NewMainEnvironment())
			ctx = call.WithEval(ctx, func(ast.Node, env.Environment[ruby.Object]) (ruby.Object, error) {
				return object.NewInteger(8), nil
			})

			function := &object.Function{}

			result, _ := function.Call(ctx)
			assert.EqualCmpAny(t, object.NewInteger(8), result, object.CompareRubyObjectsForTests)
		})
		t.Run("wrapped into a return value", func(t *testing.T) {
			ctx := call.NewContext(context.Background(), nil, object.NewMainEnvironment())
			ctx = call.WithEval(ctx, func(ast.Node, env.Environment[ruby.Object]) (ruby.Object, error) {
				return &object.ReturnValue{Value: object.NewInteger(8)}, nil
			})

			function := &object.Function{}

			result, _ := function.Call(ctx)
			assert.EqualCmpAny(t, object.NewInteger(8), result, object.CompareRubyObjectsForTests)
		})
	})
	t.Run("validates that the arguments match the function parameters", func(t *testing.T) {
		ctx := call.NewContext(context.Background(), nil, object.NewMainEnvironment())
		ctx = call.WithEval(ctx, func(ast.Node, env.Environment[ruby.Object]) (ruby.Object, error) {
			return nil, nil
		})

		function := &object.Function{
			Parameters: []*object.FunctionParameter{},
		}

		t.Run("without block argument", func(t *testing.T) {
			_, err := function.Call(ctx, object.NewString("foo"))
			assert.Error(t, err, object.NewWrongNumberOfArgumentsError(0, 1))
		})

		t.Run("with default arguments", func(t *testing.T) {
			function.Parameters = []*object.FunctionParameter{
				{Name: "x", Default: object.TRUE},
				{Name: "y"},
			}

			_, err := function.Call(ctx, object.NewInteger(8))
			assert.NoError(t, err)
		})
	})
}
