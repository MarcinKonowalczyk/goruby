package object

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/pkg/errors"
)

func TestObjectIsNil(t *testing.T) {
	context := &callContext{receiver: TRUE}
	result, err := objectIsNil(context)

	checkError(t, err, nil)

	boolean, ok := SymbolToBool(result)
	if !ok {
		t.Logf("Expected Boolean, got %T", result)
		t.FailNow()
	}

	if boolean != false {
		t.Logf("Expected false, got true")
		t.Fail()
	}
}

func TestObjectRequire(t *testing.T) {
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
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		result, err := objectRequire(context, name)

		if err != nil {
			t.Logf("expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := SymbolToBool(result)
		if !ok {
			t.Logf("Expected Boolean, got %#v", result)
			t.FailNow()
		}

		if result != TRUE {
			t.Logf("Expected return to equal TRUE, got FALSE")
			t.Fail()
		}

		if evalCallCount != 1 {
			t.Logf("Expected context.Eval to be called once, was %d\n", evalCallCount)
			t.Fail()
		}

		expectedASTNodeString := "x = 5"
		actualASTNodeString := evalCallASTNode.String()
		if expectedASTNodeString != actualASTNodeString {
			t.Logf("Expected Eval AST param to equal %q, got %q\n", expectedASTNodeString, actualASTNodeString)
			t.Fail()
		}
	})
	t.Run("env side effects no $LOADED_FEATURES", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		_, err := objectRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("env side effects missing suffix", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := objectRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("env side effects $LOADED_FEATURES exist", func(t *testing.T) {
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(&String{"foo"}))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := objectRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{"foo"}, &String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
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
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := objectRequire(context, name)
		if err != nil {
			panic(err)
		}

		_, ok := env.Get("x")

		if ok {
			t.Logf("Expected local variable not to leak over require")
			t.Fail()
		}
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
			receiver: &Object{},
		}
		name := &String{"file/not/exist"}

		_, err := objectRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		expectedErr := NewNoSuchFileLoadError("file/not/exist")
		if !reflect.DeepEqual(expectedErr, err) {
			t.Logf("Expected error to equal %v, got %v", expectedErr, err)
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("syntax error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile_syntax_error.rb"}

		_, err := objectRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		syntaxErr, ok := err.(*SyntaxError)
		if !ok {
			t.Logf("Expected syntax error, got %T:%v\n", err, err)
			t.Fail()
		}
		underlyingErr := syntaxErr.UnderlyingError()
		if !parser.IsEOFError(underlyingErr) {
			t.Logf("Expected EOF error, got:\n%q", underlyingErr)
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("thrown error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return nil, NewException("something went wrong")
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile_name_error.rb"}

		_, err := objectRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		expectedErr := NewException("something went wrong")
		if !reflect.DeepEqual(expectedErr, errors.Cause(err)) {
			t.Logf("Expected error to equal\n%q\n\tgot\n%q", expectedErr, err)
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("already loaded", func(t *testing.T) {
		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(&String{abs}))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		result, err := objectRequire(context, name)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.Fail()
		}

		_, ok := SymbolToBool(result)
		if !ok {
			t.Logf("Expected Boolean, got %#v", result)
			t.FailNow()
		}

		if result != FALSE {
			t.Logf("Expected return to equal FALSE, got TRUE")
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
}

func TestObjectToS(t *testing.T) {
	t.Run("object as receiver", func(t *testing.T) {
		context := &callContext{
			receiver: &Object{},
		}

		result, err := objectToS(context)

		checkError(t, err, nil)

		expected := &String{Value: fmt.Sprintf("#<Object:%p>", context.receiver)}

		checkResult(t, result, expected)
	})
	t.Run("self object as receiver", func(t *testing.T) {
		self := &Self{RubyObject: &Object{}, Name: "foo"}
		context := &callContext{
			receiver: self,
		}

		result, err := objectToS(context)

		checkError(t, err, nil)

		expected := &String{Value: fmt.Sprintf("#<Object:%p>", self.RubyObject)}

		checkResult(t, result, expected)
	})
}

func TestObjectRaise(t *testing.T) {
	object := &Self{RubyObject: &Object{}, Name: "x"}
	env := NewMainEnvironment()
	context := &callContext{
		receiver: object,
		env:      env,
	}

	t.Run("without args", func(t *testing.T) {
		result, err := objectRaise(context)

		checkResult(t, result, nil)

		checkError(t, err, NewRuntimeError(""))
	})

	t.Run("with 1 arg", func(t *testing.T) {
		t.Run("string argument", func(t *testing.T) {
			result, err := objectRaise(context, &String{Value: "ouch"})

			checkResult(t, result, nil)

			checkError(t, err, NewRuntimeError("ouch"))
		})
		t.Run("class argument", func(t *testing.T) {
			t.Run("exception class", func(t *testing.T) {
				result, err := objectRaise(context, standardErrorClass)

				checkResult(t, result, nil)

				checkError(t, err, &StandardError{message: "StandardError"})
			})
			t.Run("other class", func(t *testing.T) {
				result, err := objectRaise(context, stringClass)

				checkResult(t, result, nil)

				checkError(t, err, &TypeError{Message: "exception class/object expected"})
			})
			t.Run("object with #exception returning exception", func(t *testing.T) {
				exceptionFn := func(CallContext, ...RubyObject) (RubyObject, error) {
					return &StandardError{message: "err"}, nil
				}
				obj := &extendedObject{
					RubyObject: &Object{},
					class: newEigenclass(objectClass, map[string]RubyMethod{
						"exception": publicMethod(exceptionFn),
					}),
					Environment: NewEnvironment(),
				}

				result, err := objectRaise(context, obj)

				checkResult(t, result, nil)

				checkError(t, err, &StandardError{message: "err"})
			})
		})
	})
}
