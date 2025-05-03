package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestExceptionInitialize(t *testing.T) {
	context := &callContext{
		receiver: &Exception{},
		env:      NewMainEnvironment(),
	}
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionInitialize(context)

		utils.AssertNoError(t, err)

		utils.AssertEqualCmpAny(t, result, &Exception{message: "Exception"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			result, err := exceptionInitialize(context, NewString("err"))

			utils.AssertNoError(t, err)

			utils.AssertEqualCmpAny(t, result, &Exception{message: "err"}, CompareRubyObjectsForTests)
		})
		t.Run("other object", func(t *testing.T) {
			result, err := exceptionInitialize(context, &Symbol{Value: "symbol"})

			utils.AssertNoError(t, err)

			utils.AssertEqualCmpAny(t, result, &Exception{message: "symbol"}, CompareRubyObjectsForTests)
		})
	})
}

func TestExceptionException(t *testing.T) {
	contextObject := &Exception{message: "x"}
	context := &callContext{
		receiver: contextObject,
		env:      NewMainEnvironment(),
	}
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionException(context)

		utils.AssertNoError(t, err)

		if contextObject != result {
			t.Logf("Expected result to pointer equal context\n")
			t.Fail()
		}
		utils.AssertEqualCmpAny(t, result, &Exception{message: "x"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg", func(t *testing.T) {
		result, err := exceptionException(context, NewString("x"))

		utils.AssertNoError(t, err)

		if contextObject != result {
			t.Logf("Expected result to pointer equal context\n")
			t.Fail()
		}

		utils.AssertEqualCmpAny(t, result, &Exception{message: "x"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg but different message", func(t *testing.T) {
		result, err := exceptionException(context, NewString("err"))

		utils.AssertNoError(t, err)

		utils.AssertEqualCmpAny(t, result, &Exception{message: "err"}, CompareRubyObjectsForTests)
	})
}

func TestExceptionToS(t *testing.T) {
	contextObject := &Exception{message: "x"}
	context := &callContext{
		receiver: contextObject,
		env:      NewMainEnvironment(),
	}

	result, err := exceptionToS(context)

	utils.AssertNoError(t, err)

	utils.AssertEqualCmpAny(t, result, NewString("x"), CompareRubyObjectsForTests)
}
