package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestInteger_hashKey(t *testing.T) {
	hello1 := &Integer{Value: 1}
	hello2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 3}
	diff2 := &Integer{Value: 3}

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

func TestIntegerDiv(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(2),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), &Integer{}),
		},
		{
			[]RubyObject{NewInteger(0)},
			nil,
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerDiv(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerMul(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(8),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerMul(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)
		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(4),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(2)}

		result, err := integerAdd(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)
		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerSub(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(3)},
			NewInteger(1),
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewCoercionTypeError(NewString(""), &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerSub(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestIntegerModulo(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewInteger(3)},
// 			NewInteger(1),
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewCoercionTypeError(NewString(""), &Integer{}),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewInteger(4)}

// 		result, err := integerModulo(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

func TestIntegerLt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerLt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerGt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerGt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

// func TestIntegerEq(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewInteger(6)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewInteger(4)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Integer with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewInteger(4)}

// 		result, err := integerEq(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

// func TestIntegerNeq(t *testing.T) {
// 	tests := []struct {
// 		arguments []RubyObject
// 		result    RubyObject
// 		err       error
// 	}{
// 		{
// 			[]RubyObject{NewInteger(6)},
// 			TRUE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewInteger(4)},
// 			FALSE,
// 			nil,
// 		},
// 		{
// 			[]RubyObject{NewString("")},
// 			nil,
// 			NewArgumentError("comparison of Integer with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewInteger(4)}

// 		result, err := integerNeq(context, testCase.arguments...)

// 		utils.AssertError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
// 	}
// }

func TestIntegerGte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerGte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerLte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewString("")},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerLte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}

func TestIntegerSpaceship(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			&Integer{Value: -1},
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			&Integer{Value: 0},
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			&Integer{Value: 1},
			nil,
		},
		{
			[]RubyObject{NewString("")},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerSpaceship(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualCmpAny(t, result, testCase.result, CompareRubyObjectsForTests)
	}
}
