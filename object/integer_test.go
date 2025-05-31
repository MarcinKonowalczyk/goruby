package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestInteger_hashKey(t *testing.T) {
	hello1 := NewInteger(1)
	hello2 := NewInteger(1)
	diff1 := NewInteger(3)
	diff2 := NewInteger(3)

	assert.Equal(t, hello1.HashKey(), hello2.HashKey())
	assert.Equal(t, diff1.HashKey(), diff2.HashKey())
	assert.NotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func TestIntegerDiv(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(2)},
			NewInteger(2),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewInteger(0)),
		},
		{
			[]ruby.Object{NewInteger(0)},
			nil,
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, err := integerDiv(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerMul(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(2)},
			NewInteger(8),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewInteger(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, err := integerMul(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerAdd(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(2)},
			NewInteger(4),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewInteger(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(2), nil)
		result, err := integerAdd(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerSub(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(3)},
			NewInteger(1),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewInteger(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, err := integerSub(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestIntegerModulo(t *testing.T) {
// 	tests := []struct {
// 		arguments []ruby.Object
// 		result    ruby.Object
// 		err       error
// 	}{
// 		{
// 			[]ruby.Object{NewInteger(3)},
// 			NewInteger(1),
// 			nil,
// 		},
// 		{
// 			[]ruby.Object{NewString("")},
// 			nil,
// 			NewCoercionTypeError(NewString(""), NewInteger(0)),)),
// 		},
// 	}

// 	for _, testCase := range tests {
// ctx := call.NewContext[ruby.Object](NewInteger(4), nil)

// 		result, err := integerModulo(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

func TestIntegerLt(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, _ := integerLt(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerGt(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, _ := integerGt(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestIntegerEq(t *testing.T) {
// 	tests := []struct {
// 		arguments []ruby.Object
// 		result    ruby.Object
// 		err       error
// 	}{
// 		{
// 			[]ruby.Object{NewInteger(6)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]ruby.Object{NewInteger(4)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]ruby.Object{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Integer with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// ctx := call.NewContext[ruby.Object](NewInteger(4), nil)

// 		result, err := integerEq(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

// func TestIntegerNeq(t *testing.T) {
// 	tests := []struct {
// 		arguments []ruby.Object
// 		result    ruby.Object
// 		err       error
// 	}{
// 		{
// 			[]ruby.Object{NewInteger(6)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]ruby.Object{NewInteger(4)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]ruby.Object{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Integer with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// ctx := call.NewContext[ruby.Object](NewInteger(4), nil)

// 		result, err := integerNeq(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

func TestIntegerGte(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, _ := integerGte(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerLte(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, _ := integerLte(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerSpaceship(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewInteger(6)},
			NewInteger(-1),
			nil,
		},
		{
			[]ruby.Object{NewInteger(4)},
			NewInteger(0),
			nil,
		},
		{
			[]ruby.Object{NewInteger(2)},
			NewInteger(1),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](NewInteger(4), nil)
		result, _ := integerSpaceship(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
