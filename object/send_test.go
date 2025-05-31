package object

import (
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
	"github.com/pkg/errors"
)

type testRubyObject struct {
	class ruby.ClassObject
	Name  string
}

func (t *testRubyObject) Inspect() string { return "TEST OBJECT" }
func (t *testRubyObject) Class() ruby.Class {
	if t.class != nil {
		return t.class
	}
	return bottomClass
}
func (t *testRubyObject) HashKey() hash.Key {
	return hash.Key(99)
}

func TestSend(t *testing.T) {
	// }
	methods := map[string]ruby.Method{
		"a_method": newMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return TRUE, nil
		}),
		"another_method": newMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return FALSE, nil
		}),
	}
	t.Run("normal object as context", func(t *testing.T) {
		// ctx := &callContext{
		// 	receiver: &testRubyObject{
		// 		class: &class{
		// 			name:            "base class",
		// 			instanceMethods: NewMethodSet(methods),
		// 		},
		// 	},
		// }
		ctx := call.NewContext[ruby.Object](
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(methods),
				},
			},
			nil,
		)

		tests := []struct {
			method         string
			expectedResult ruby.Object
			expectedError  error
		}{
			{
				"a_method",
				TRUE,
				nil,
			},
			{
				"another_method",
				FALSE,
				nil,
			},
			{
				"unknown_method",
				nil,
				NewNoMethodError(ctx.Receiver(), "unknown_method"),
			},
		}

		for _, testCase := range tests {
			result, err := Send(ctx, testCase.method)

			assert.Error(t, errors.Cause(err), testCase.expectedError)
			assert.EqualCmpAny(t, result, testCase.expectedResult, CompareRubyObjectsForTests)
		}
	})
}

func TestAddMethod(t *testing.T) {
	t.Run("vanilla object", func(t *testing.T) {
		call_context := &testRubyObject{
			class: &class{
				name:            "base class",
				instanceMethods: NewMethodSet(map[string]ruby.Method{}),
			},
		}

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  env.NewEnvironment[ruby.Object](),
			Body: nil,
		}

		newContext, _ := AddMethod(call_context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		assert.That(t, ok, "Expected object to have method foo")
	})
	t.Run("class object", func(t *testing.T) {
		call_context := &class{
			name:            "A",
			instanceMethods: NewMethodSet(map[string]ruby.Method{}),
		}

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  env.NewEnvironment[ruby.Object](),
			Body: nil,
		}

		newContext, _ := AddMethod(call_context, "foo", fn)

		class, ok := newContext.(*class)
		assert.That(t, ok, "Expected returned object to be a class, got %T", newContext)

		_, ok = class.Methods().Get("foo")
		assert.That(t, ok, "Expected object to have method foo")
	})
	t.Run("extended object", func(t *testing.T) {
		call_context := newExtendedObject(
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(map[string]ruby.Method{}),
				},
			},
		)

		call_context.eigenclass.addMethod("bar", newMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return NIL, nil
		}))

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  env.NewEnvironment[ruby.Object](),
			Body: nil,
		}

		newContext, _ := AddMethod(call_context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		assert.That(t, ok, "Expected object to have method foo")

		_, ok = newContext.Class().Methods().Get("bar")
		assert.That(t, ok, "Expected object to have method bar")
	})
	t.Run("vanilla self object", func(t *testing.T) {
		vanillaObject := &testRubyObject{
			class: &class{
				name:            "base class",
				instanceMethods: NewMethodSet(map[string]ruby.Method{}),
			},
			Name: "main",
		}
		call_context := vanillaObject

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  env.NewEnvironment[ruby.Object](),
			Body: nil,
		}

		newContext, extended := AddMethod(call_context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		assert.That(t, ok, "Expected object to have method foo")
		assert.That(t, extended, "Expected object to be extended")

		returnPointer := reflect.ValueOf(newContext).Pointer()
		contextPointer := reflect.ValueOf(call_context).Pointer()
		assert.NotEqual(t, returnPointer, contextPointer)

		extendedRubyObject := newContext.(*extendedObject).Object
		assert.EqualCmpAny(t, vanillaObject, extendedRubyObject, CompareRubyObjectsForTests)
	})
	t.Run("extended self object", func(t *testing.T) {
		call_context := newExtendedObject(
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(map[string]ruby.Method{}),
				},
			},
		)

		call_context.eigenclass.addMethod("bar", newMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return NIL, nil
		}))

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  env.NewEnvironment[ruby.Object](),
			Body: nil,
		}

		newContext, extended := AddMethod(call_context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		assert.That(t, ok, "Expected object to have method foo")

		_, ok = newContext.Class().Methods().Get("bar")
		assert.That(t, ok, "Expected object to have method bar")

		returnPointer := reflect.ValueOf(newContext).Pointer()
		contextPointer := reflect.ValueOf(call_context).Pointer()

		assert.That(t, !extended, "Expected object not to be extended")
		assert.Equal(t, returnPointer, contextPointer)
	})
}
