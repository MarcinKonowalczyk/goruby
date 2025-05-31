package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestString_hashKey(t *testing.T) {
	hello1 := NewString("Hello World")
	hello2 := NewString("Hello World")
	diff1 := NewString("My name is johnny")
	diff2 := NewString("My name is johnny")

	assert.Equal(t, hello1.HashKey(), hello2.HashKey())
	assert.Equal(t, diff1.HashKey(), diff2.HashKey())
	assert.NotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func Test_stringify(t *testing.T) {
	t.Run("object with regular `to_s`", func(t *testing.T) {
		ctx := call.NewContext[ruby.Object](nil, nil)
		res, err := stringify(ctx, NewSymbol("sym"))
		assert.NoError(t, err)
		assert.Equal(t, res, "sym")
	})
	t.Run("object without `to_s`", func(t *testing.T) {
		ctx := call.NewContext[ruby.Object](nil, nil)
		_, err := stringify(ctx, nil)
		assert.Error(t, err, NewTypeError("can't convert nil into String"))
	})
}

func TestStringAdd(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewString(" bar")},
			NewString("foo bar"),
			nil,
		},
		{
			[]ruby.Object{NewInteger(3)},
			nil,
			NewImplicitConversionTypeError(NewString(""), NewInteger(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewString("foo"), nil)
		result, err := stringAdd(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func Test_StringGsub(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewString("o"), NewString("zz")},
			NewString("fzzzzbar"),
			nil,
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewString("foobar"), nil)
		result, err := stringGsub(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
