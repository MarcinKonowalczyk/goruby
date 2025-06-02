package object

import (
	"context"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestBottomIsNil(t *testing.T) {
	ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(TRUE)
	result, err := bottomIsNil(ctx)

	assert.NoError(t, err)

	boolean, ok := SymbolToBool(result)
	assert.That(t, ok, "Expected boolean, got %T", result)
	assert.That(t, !boolean, "Expected false, got true")
}

func TestBottomToS(t *testing.T) {
	t.Run("object as receiver", func(t *testing.T) {
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(&Bottom{})
		result, err := bottomToS(ctx)
		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewStringf("#<Bottom:%p>", ctx.Receiver()), CompareRubyObjectsForTests)
	})
	t.Run("self object as receiver", func(t *testing.T) {
		self := &Bottom{}
		ctx := call.NewContext[ruby.Object](context.Background(), nil).WithReceiver(self)
		result, err := bottomToS(ctx)
		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewStringf("#<Bottom:%p>", self), CompareRubyObjectsForTests)
	})
}

func TestBottomRaise(t *testing.T) {
	object := &Bottom{}
	env := NewMainEnvironment()
	ctx := call.NewContext[ruby.Object](context.Background(), nil).
		WithReceiver(object).
		WithEnv(env)

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
