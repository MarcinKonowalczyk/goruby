package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestFloat_hashKey(t *testing.T) {
	hello1 := NewFloat(1)
	hello2 := NewFloat(1)
	diff1 := NewFloat(3)
	diff2 := NewFloat(3)
	utils.AssertEqual(t, hello1.HashKey(), hello2.HashKey())
	utils.AssertEqual(t, diff1.HashKey(), diff2.HashKey())
	utils.AssertNotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func TestFloatDiv(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(2)},
			NewFloat(2),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
		{
			[]RubyObject{NewFloat(0)},
			nil,
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatDiv(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatMul(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(2)},
			NewFloat(8),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatMul(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(2)},
			NewFloat(4),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(2)}

		result, err := floatAdd(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatSub(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(3)},
			NewFloat(1),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), NewFloat(0)),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatSub(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
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
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatModulo(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
// 	}
// }

func TestFloatLt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, _ := floatLt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatGt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, _ := floatGt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
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
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatEq(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
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
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatNeq(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjects)
// 	}
// }

func TestFloatGte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewFloat(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, _ := floatGte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatLte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewFloat(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, _ := floatLte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestFloatSpaceship(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewFloat(6)},
			NewFloat(-1),
			nil,
		},
		{
			[]RubyObject{NewFloat(4)},
			NewFloat(0),
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			NewFloat(1),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, _ := floatSpaceship(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
