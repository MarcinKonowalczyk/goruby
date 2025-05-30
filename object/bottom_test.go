package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestBottomIsNil(t *testing.T) {
	ctx := &callContext{receiver: TRUE}
	result, err := bottomIsNil(ctx)

	assert.NoError(t, err)

	boolean, ok := SymbolToBool(result)
	assert.That(t, ok, "Expected boolean, got %T", result)
	assert.That(t, !boolean, "Expected false, got true")
}

func TestBottomToS(t *testing.T) {
	t.Run("object as receiver", func(t *testing.T) {
		ctx := &callContext{receiver: &Bottom{}}
		result, err := bottomToS(ctx)
		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewStringf("#<Bottom:%p>", ctx.receiver), CompareRubyObjectsForTests)
	})
	t.Run("self object as receiver", func(t *testing.T) {
		self := &Bottom{}
		ctx := &callContext{receiver: self}
		result, err := bottomToS(ctx)
		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewStringf("#<Bottom:%p>", self), CompareRubyObjectsForTests)
	})
}

func TestBottomRaise(t *testing.T) {
	object := &Bottom{}
	env := NewMainEnvironment()
	ctx := &callContext{
		receiver: object,
		env:      env,
	}

	t.Run("without args", func(t *testing.T) {
		result, err := bottomRaise(ctx)
		assert.EqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
		assert.Error(t, err, NewRuntimeError(""))
	})

	t.Run("with 1 arg", func(t *testing.T) {
		t.Run("string argument", func(t *testing.T) {
			result, err := bottomRaise(ctx, NewString("ouch"))
			assert.EqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			assert.Error(t, err, NewRuntimeError("ouch"))
		})
		t.Run("integer argument", func(t *testing.T) {
			obj := NewInteger(5)
			result, err := bottomRaise(ctx, obj)
			assert.EqualCmpAny(t, result, nil, CompareRubyObjectsForTests)
			assert.Error(t, err, NewRuntimeError("%s", obj.Inspect()))
		})
	})
}
