package object

import (
	"testing"
)

func TestFloat_hashKey(t *testing.T) {
	hello1 := &Float{Value: 1}
	hello2 := &Float{Value: 1}
	diff1 := &Float{Value: 3}
	diff2 := &Float{Value: 3}

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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Float{}),
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

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Float{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatMul(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Float{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(2)}

		result, err := floatAdd(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Float{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatSub(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
// 			[]RubyObject{&String{""}},
// 			nil,
// 			NewCoercionTypeError(&String{}, &Float{}),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatModulo(context, testCase.arguments...)

// 		checkError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatLt(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatGt(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
// 			[]RubyObject{&String{""}},
// 			nil,
// 			NewArgumentError("comparison of Float with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatEq(context, testCase.arguments...)

// 		checkError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
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
// 			[]RubyObject{&String{""}},
// 			nil,
// 			NewArgumentError("comparison of Float with String failed"),
// 		},
// 	}

// 	for _, testCase := range tests {
// 		context := &callContext{receiver: NewFloat(4)}

// 		result, err := floatNeq(context, testCase.arguments...)

// 		checkError(t, err, testCase.err)

// 		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatGte(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Float with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatLte(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			&Float{Value: -1},
			nil,
		},
		{
			[]RubyObject{NewFloat(4)},
			&Float{Value: 0},
			nil,
		},
		{
			[]RubyObject{NewFloat(2)},
			&Float{Value: 1},
			nil,
		},
		{
			[]RubyObject{&String{""}},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewFloat(4)}

		result, err := floatSpaceship(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}
