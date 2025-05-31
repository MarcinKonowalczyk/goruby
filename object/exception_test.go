package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestExceptionInitialize(t *testing.T) {
	ctx := call.NewContext[ruby.Object](&Exception{}, NewMainEnvironment())
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
