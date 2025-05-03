package object

import (
	"path/filepath"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/MarcinKonowalczyk/goruby/utils"
	"github.com/pkg/errors"
)

func TestBottomIsNil(t *testing.T) {
	context := &callContext{receiver: TRUE}
	result, err := bottomIsNil(context)

	utils.AssertNoError(t, err)

	boolean, ok := SymbolToBool(result)
	utils.Assert(t, ok, "Expected boolean, got %T", result)
	utils.Assert(t, !boolean, "Expected false, got true")
}

func TestBottomRequire(t *testing.T) {
	t.Run("wiring together", func(t *testing.T) {
		evalCallCount := 0
		var evalCallASTNode ast.Node
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			evalCallCount++
			evalCallASTNode = node
			return TRUE, nil
		}

		context := &callContext{
			env:      NewEnvironment(),
			eval:     eval,
			receiver: &Bottom{},
		}
		name := NewString("./fixtures/testfile.rb")

		result, err := bottomRequire(context, name)
		utils.AssertNoError(t, err)

		_, ok := SymbolToBool(result)
		utils.Assert(t, ok, "Expected boolean, got %T", result)
		utils.AssertEqualCmpAny(t, result, TRUE, CompareRubyObjectsForTests)
		utils.AssertEqual(t, evalCallCount, 1)
		utils.AssertEqual(t, "x = 5", evalCallASTNode.String())
	})
	t.Run("env side effects no $LOADED_FEATURES", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		name := NewString("./fixtures/testfile.rb")

		_, err := bottomRequire(context, name)
		utils.AssertNoError(t, err)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		utils.AssertEqualCmpAny(t, arr, NewArray(NewString(abs)), CompareRubyObjectsForTests)
	})
	t.Run("env side effects missing suffix", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		name := NewString("./fixtures/testfile")

		_, err := bottomRequire(context, name)
		utils.AssertNoError(t, err)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		utils.AssertEqualCmpAny(t, arr, NewArray(NewString(abs)), CompareRubyObjectsForTests)
	})
	t.Run("env side effects $LOADED_FEATURES exist", func(t *testing.T) {
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(NewString("foo")))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		name := NewString("./fixtures/testfile")

		_, err := bottomRequire(context, name)
		utils.AssertNoError(t, err)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		utils.AssertEqualCmpAny(t, arr, NewArray(NewString("foo"), NewString(abs)), CompareRubyObjectsForTests)
	})
	t.Run("env side effects local variables", func(t *testing.T) {
		env := NewEnvironment()
		var eval func(node ast.Node, env Environment) (RubyObject, error)
		eval = func(node ast.Node, env Environment) (RubyObject, error) {
			switch node := node.(type) {
			case *ast.Program:
				var result RubyObject
				var err error
				for _, statement := range node.Statements {
					result, err = eval(statement, env)

					if err != nil {
						return nil, err
					}
				}
				return result, nil
			case *ast.ExpressionStatement:
				return eval(node.Expression, env)
			case *ast.Assignment:
				val, err := eval(node.Right, env)
				if err != nil {
					return nil, err
				}
				env.Set(node.String(), val)
				return val, nil
			}
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		_, err := bottomRequire(context, NewString("./fixtures/testfile"))
		utils.AssertNoError(t, err)

		_, ok := env.Get("x")
		utils.Assert(t, !ok, "Expected local variable not to leak over require")
	})
	t.Run("file does not exist", func(t *testing.T) {
		env := NewEnvironment()
		env.SetGlobal("$:", NewArray())
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		_, err := bottomRequire(context, NewString("file/not/exist"))
		utils.AssertError(t, err, NewNoSuchFileLoadError("file/not/exist"))

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		utils.AssertEqualCmpAny(t, arr, NewArray(), CompareRubyObjectsForTests)
	})
	t.Run("syntax error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		_, err := bottomRequire(context, NewString("./fixtures/testfile_syntax_error.rb"))
		utils.AssertNotEqual(t, err, nil)

		syntaxErr, ok := err.(*SyntaxError)
		utils.Assert(t, ok, "Expected SyntaxError, got %T:%v", err, err)

		underlyingErr := syntaxErr.UnderlyingError()
		utils.Assert(t, parser.IsEOFError(underlyingErr), "Expected EOF error, got:\n%q", underlyingErr)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
		utils.AssertEqualCmpAny(t, arr, NewArray(), CompareRubyObjectsForTests)
	})
	t.Run("thrown error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return nil, NewException("something went wrong")
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		_, err := bottomRequire(context, NewString("./fixtures/testfile_name_error.rb"))
		utils.AssertNotEqual(t, err, nil)

		expectedErr := NewException("something went wrong")
		utils.AssertEqualCmpAny(t, errors.Cause(err), expectedErr, CompareRubyObjectsForTests)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		expected := NewArray()
		utils.AssertEqualCmpAny(t, arr, expected, CompareRubyObjectsForTests)
	})
	t.Run("already loaded", func(t *testing.T) {
		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(NewString(abs)))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Bottom{},
		}
		name := NewString("./fixtures/testfile.rb")

		result, err := bottomRequire(context, name)
		utils.AssertNoError(t, err)

		_, ok := SymbolToBool(result)
		utils.Assert(t, ok, "Expected boolean, got %T", result)
		utils.AssertEqualCmpAny(t, result, FALSE, CompareRubyObjectsForTests)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		utils.Assert(t, ok, "Expected env to contain global $LOADED_FEATURES")

		arr, ok := loadedFeatures.(*Array)
		utils.Assert(t, ok, "Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)

		expected := NewArray(NewString(abs))
		utils.AssertEqualCmpAny(t, arr, expected, CompareRubyObjectsForTests)
	})
}

func TestBottomToS(t *testing.T) {
	t.Run("object as receiver", func(t *testing.T) {
		context := &callContext{receiver: &Bottom{}}
		result, err := bottomToS(context)
		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewStringf("#<Bottom:%p>", context.receiver), CompareRubyObjectsForTests)
	})
	t.Run("self object as receiver", func(t *testing.T) {
		self := &Bottom{}
		context := &callContext{receiver: self}
		result, err := bottomToS(context)
		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewStringf("#<Bottom:%p>", self), CompareRubyObjectsForTests)
	})
}

func TestBottomRaise(t *testing.T) {
	object := &Bottom{}
	env := NewMainEnvironment()
	context := &callContext{
		receiver: object,
		env:      env,
	}

	t.Run("without args", func(t *testing.T) {
		result, err := bottomRaise(context)
		utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
		utils.AssertError(t, err, NewRuntimeError(""))
	})

	t.Run("with 1 arg", func(t *testing.T) {
		t.Run("string argument", func(t *testing.T) {
			result, err := bottomRaise(context, NewString("ouch"))
			utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			utils.AssertError(t, err, NewRuntimeError("ouch"))
		})
		t.Run("integer argument", func(t *testing.T) {
			obj := &Integer{Value: 5}
			result, err := bottomRaise(context, obj)
			utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			utils.AssertError(t, err, NewRuntimeError("%s", obj.Inspect()))
		})
	})
}
