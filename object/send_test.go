package object

import (
	"reflect"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
	"github.com/pkg/errors"
)

type testRubyObject struct {
	class RubyClassObject
	Name  string
}

func (t *testRubyObject) Inspect() string { return "TEST OBJECT" }
func (t *testRubyObject) Class() RubyClass {
	if t.class != nil {
		return t.class
	}
	return bottomClass
}
func (t *testRubyObject) HashKey() HashKey {
	return HashKey(99)
}

func TestSend(t *testing.T) {
	// }
	methods := map[string]RubyMethod{
		"a_method": newMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return TRUE, nil
		}),
		"another_method": newMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return FALSE, nil
		}),
	}
	t.Run("normal object as context", func(t *testing.T) {
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(methods),
				},
			},
		}

		tests := []struct {
			method         string
			expectedResult RubyObject
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
				NewNoMethodError(context.receiver, "unknown_method"),
			},
		}

		for _, testCase := range tests {
			result, err := Send(context, testCase.method)

			utils.AssertError(t, errors.Cause(err), testCase.expectedError)
			utils.AssertEqualCmpAny(t, result, testCase.expectedResult, CompareRubyObjectsForTests)
		}
	})
}

func TestAddMethod(t *testing.T) {
	t.Run("vanilla object", func(t *testing.T) {
		context := &testRubyObject{
			class: &class{
				name:            "base class",
				instanceMethods: NewMethodSet(map[string]RubyMethod{}),
			},
		}

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext, _ := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}
	})
	t.Run("class object", func(t *testing.T) {
		context := &class{
			name:            "A",
			instanceMethods: NewMethodSet(map[string]RubyMethod{}),
		}

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext, _ := AddMethod(context, "foo", fn)

		class, ok := newContext.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", newContext)
			t.FailNow()
		}

		_, ok = class.Methods().Get("foo")
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}
	})
	t.Run("extended object", func(t *testing.T) {
		context := newExtendedObject(
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(map[string]RubyMethod{}),
				},
			},
		)

		context.eigenclass.addMethod("bar", newMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return NIL, nil
		}))

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext, _ := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		_, ok = newContext.Class().Methods().Get("bar")
		if !ok {
			t.Logf("Expected object to have method bar")
			t.Fail()
		}
	})
	t.Run("vanilla self object", func(t *testing.T) {
		vanillaObject := &testRubyObject{
			class: &class{
				name:            "base class",
				instanceMethods: NewMethodSet(map[string]RubyMethod{}),
			},
			Name: "main",
		}
		context := vanillaObject

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext, extended := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		returnPointer := reflect.ValueOf(newContext).Pointer()
		contextPointer := reflect.ValueOf(context).Pointer()

		if !extended {
			t.Logf("Expected object to be extended")
			t.Fail()
		}

		if returnPointer == contextPointer {
			t.Logf("Expected input and return context to be different")
			t.Fail()
		}

		extendedRubyObject := newContext.(*extendedObject).RubyObject

		if !reflect.DeepEqual(vanillaObject, extendedRubyObject) {
			t.Logf(
				"Expected wrapped extended object to equal\n%+#v\n\tgot\n%+#v\n",
				vanillaObject,
				extendedRubyObject,
			)
			t.Fail()
		}
	})
	t.Run("extended self object", func(t *testing.T) {
		context := newExtendedObject(
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: NewMethodSet(map[string]RubyMethod{}),
				},
			},
		)

		context.eigenclass.addMethod("bar", newMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return NIL, nil
		}))

		fn := &Function{
			Parameters: []*FunctionParameter{
				{Name: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext, extended := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		_, ok = newContext.Class().Methods().Get("bar")
		if !ok {
			t.Logf("Expected object to have method bar")
			t.Fail()
		}

		returnPointer := reflect.ValueOf(newContext).Pointer()
		contextPointer := reflect.ValueOf(context).Pointer()

		if extended {
			t.Logf("Expected object not to be extended")
			t.Fail()
		}

		if returnPointer != contextPointer {
			t.Logf("Expected input and return context to be the same")
			t.Fail()
		}
	})
}
