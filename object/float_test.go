package object

import (
	"context"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestFloat_hashKey(t *testing.T) {
	hello1 := NewFloat(1)
	hello2 := NewFloat(1)
	diff1 := NewFloat(3)
	diff2 := NewFloat(3)
	assert.Equal(t, hello1.HashKey(), hello2.HashKey())
	assert.Equal(t, diff1.HashKey(), diff2.HashKey())
	assert.NotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func TestFloatDiv(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(2)},
			NewFloat(2),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
		{
			[]ruby.Object{NewFloat(0)},
			nil,
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, err := floatDiv(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatMul(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(2)},
			NewFloat(8),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, err := floatMul(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatAdd(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(2)},
			NewFloat(4),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(2))
		result, err := floatAdd(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatSub(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(3)},
			NewFloat(1),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, err := floatSub(ctx, testCase.arguments...)
		assert.Error(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestFloatModulo(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewFloat(3)},
// 			NewFloat(1),
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewCoercionTypeError(NewString(""), NewFloat(0)),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))

// 		result, err := floatModulo(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		assert.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
// 	}
// }

func TestFloatLt(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(6)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(2)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, _ := floatLt(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatGt(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(6)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(2)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, _ := floatGt(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestFloatEq(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewFloat(6)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewFloat(4)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Float with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))

// 		result, err := floatEq(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		assert.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
// 	}
// }

// func TestFloatNeq(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewFloat(6)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewFloat(4)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Float with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))

// 		result, err := floatNeq(ctx, testCase.arguments...)

// 		assert.AssertError(t, err, testCase.err)

// 		assert.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
// 	}
// }

func TestFloatGte(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(6)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(4)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(2)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, _ := floatGte(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatLte(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(6)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(4)},
			TRUE,
			nil,
		},
		{
			[]ruby.Object{NewFloat(2)},
			FALSE,
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, _ := floatLte(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatSpaceship(t *testing.T) {
	tests := []struct {
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			[]ruby.Object{NewFloat(6)},
			NewFloat(-1),
			nil,
		},
		{
			[]ruby.Object{NewFloat(4)},
			NewFloat(0),
			nil,
		},
		{
			[]ruby.Object{NewFloat(2)},
			NewFloat(1),
			nil,
		},
		{
			[]ruby.Object{NewString("")},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(NewFloat(4))
		result, _ := floatSpaceship(ctx, testCase.arguments...)
		// assert.AssertError(t, err, testCase.err)
		assert.EqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
