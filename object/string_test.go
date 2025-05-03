package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestString_hashKey(t *testing.T) {
	hello1 := NewString("Hello World")
	hello2 := NewString("Hello World")
	diff1 := NewString("My name is johnny")
	diff2 := NewString("My name is johnny")

	if hello1.hashKey() != hello2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.hashKey() != diff2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.hashKey() == diff1.hashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func Test_stringify(t *testing.T) {
	t.Run("object with regular `to_s`", func(t *testing.T) {
		obj := &Symbol{Value: "sym"}

		res, err := stringify(obj)

		utils.AssertNoError(t, err)

		if res != "sym" {
			t.Logf("Expected stringify to return 'sym', got %q\n", res)
			t.Fail()
		}
	})
	t.Run("object without `to_s`", func(t *testing.T) {
		_, err := stringify(nil)

		utils.AssertError(t, err, NewTypeError("can't convert nil into String"))
	})
}

func TestStringAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewString(" bar")},
			NewString("foo bar"),
			nil,
		},
		{
			[]RubyObject{NewInteger(3)},
			nil,
			NewImplicitConversionTypeError(NewString(""), NewInteger(0)),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewString("foo")}

		result, err := stringAdd(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)
		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func Test_StringGsub(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewString("o"), NewString("zz")},
			NewString("fzzzzbar"),
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewString("foobar")}

		result, err := stringGsub(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
