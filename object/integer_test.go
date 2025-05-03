package object

import (
	"fmt"
	"reflect"
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
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

		utils.AssertEqualWithComparator(t, result, testCase.result, CompareRubyObjects)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerMul(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(2)}

		result, err := integerAdd(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualWithComparator(t, result, testCase.result, CompareRubyObjects)
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
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerSub(context, testCase.arguments...)

		utils.AssertError(t, err, testCase.err)

		utils.AssertEqualWithComparator(t, result, testCase.result, CompareRubyObjects)
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
// 			[]RubyObject{&String{""}},
// 			nil,
// 			NewCoercionTypeError(&String{}, &Integer{}),
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerLt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualWithComparator(t, result, testCase.result, CompareRubyObjects)
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
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerGt(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
// 			[]RubyObject{&String{""}},
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
// 			[]RubyObject{&String{""}},
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
			[]RubyObject{&String{""}},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerGte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			NIL,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerLte(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
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
			[]RubyObject{&String{""}},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, _ := integerSpaceship(context, testCase.arguments...)

		// utils.AssertError(t, err, testCase.err)

		utils.AssertEqualWithComparator(t, result, testCase.result, CompareRubyObjects)
	}
}

func CompareRubyObjects(a, b RubyObject) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Class() != b.Class() {
		return false
	}
	if a, a_hashable := a.(hashable); a_hashable {
		if b, b_hashable := b.(hashable); b_hashable {
			fmt.Println("comparing hash keys")
			return a.hashKey() == b.hashKey()
		} else {
			// b is not hashable, we are not equal
			return false
		}
	}
	if _, b_hashable := b.(hashable); b_hashable {
		// a is not hashable, we are not equal
		return false
	}
	// ok, we are not hashable but we are the same class
	// check the addresses
	addrA := fmt.Sprintf("%p", a)
	addrB := fmt.Sprintf("%p", b)
	if addrA == addrB {
		return true
	}
	fmt.Println("comparing values")
	return reflect.DeepEqual(a, b)
}

func checkResult(t *testing.T, actual, expected RubyObject) {
	t.Helper()
	utils.AssertEqualWithComparator(t, actual, expected, CompareRubyObjects)
}
