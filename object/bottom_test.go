package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestBottomIsNil(t *testing.T) {
	context := &callContext{receiver: TRUE}
	result, err := bottomIsNil(context)

	utils.AssertNoError(t, err)

	boolean, ok := SymbolToBool(result)
	utils.Assert(t, ok, "Expected boolean, got %T", result)
	utils.Assert(t, !boolean, "Expected false, got true")
}

func TestBottomToS(t *testing.T) {
	t.Run("object as receiver", func(t *testing.T) {
		context := &callContext{receiver: &Bottom{}}
		result, err := bottomToS(context)
		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewStringf("#<Bottom:%p>", context.receiver), CompareRubyObjectsForTests)
	})
	t.Run("self object as receiver", func(t *testing.T) {
		self := &Bottom{}
		context := &callContext{receiver: self}
		result, err := bottomToS(context)
		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewStringf("#<Bottom:%p>", self), CompareRubyObjectsForTests)
	})
}

func TestBottomRaise(t *testing.T) {
	object := &Bottom{}
	env := NewMainEnvironment()
	context := &callContext{
		receiver: object,
		env:      env,
	}

	t.Run("without args", func(t *testing.T) {
		result, err := bottomRaise(context)
		utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
		utils.AssertError(t, err, NewRuntimeError(""))
	})

	t.Run("with 1 arg", func(t *testing.T) {
		t.Run("string argument", func(t *testing.T) {
			result, err := bottomRaise(context, NewString("ouch"))
			utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			utils.AssertError(t, err, NewRuntimeError("ouch"))
		})
		t.Run("integer argument", func(t *testing.T) {
			obj := NewInteger(5)
			result, err := bottomRaise(context, obj)
			utils.AssertEqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			utils.AssertError(t, err, NewRuntimeError("%s", obj.Inspect()))
		})
	})
}
