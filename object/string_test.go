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

	utils.AssertEqual(t, hello1.HashKey(), hello2.HashKey())
	utils.AssertEqual(t, diff1.HashKey(), diff2.HashKey())
	utils.AssertNotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func Test_stringify(t *testing.T) {
	t.Run("object with regular `to_s`", func(t *testing.T) {
		res, err := stringify(NewSymbol("sym"))
		utils.AssertNoError(t, err)
		utils.AssertEqual(t, res, "sym")
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

		result, err := stringAdd(context, nil, testCase.arguments...)

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

		result, err := stringGsub(context, nil, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
