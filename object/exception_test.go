package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestExceptionInitialize(t *testing.T) {
	ctx := NewCC(&Exception{}, NewMainEnvironment())
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionInitialize(ctx)

		assert.NoError(t, err)

		assert.EqualCmpAny(t, result, &Exception{message: "Exception"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			result, err := exceptionInitialize(ctx, NewString("err"))

			assert.NoError(t, err)

			assert.EqualCmpAny(t, result, &Exception{message: "err"}, CompareRubyObjectsForTests)
		})
		t.Run("other object", func(t *testing.T) {
			result, err := exceptionInitialize(ctx, NewSymbol("symbol"))

			assert.NoError(t, err)

			assert.EqualCmpAny(t, result, &Exception{message: "symbol"}, CompareRubyObjectsForTests)
		})
	})
}

func TestExceptionException(t *testing.T) {
	contextObject := &Exception{message: "x"}
	ctx := NewCC(contextObject, NewMainEnvironment())
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionException(ctx)

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, contextObject, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, result, &Exception{message: "x"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg", func(t *testing.T) {
		result, err := exceptionException(ctx, NewString("x"))

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, contextObject, CompareRubyObjectsForTests)
		assert.EqualCmpAny(t, result, &Exception{message: "x"}, CompareRubyObjectsForTests)
	})
	t.Run("with arg but different message", func(t *testing.T) {
		result, err := exceptionException(ctx, NewString("err"))

		assert.NoError(t, err)

		assert.EqualCmpAny(t, result, &Exception{message: "err"}, CompareRubyObjectsForTests)
	})
}

func TestExceptionToS(t *testing.T) {
	contextObject := &Exception{message: "x"}
	ctx := NewCC(contextObject, NewMainEnvironment())
	result, err := exceptionToS(ctx)
	assert.NoError(t, err)
	assert.EqualCmpAny(t, result, NewString("x"), CompareRubyObjectsForTests)
}
