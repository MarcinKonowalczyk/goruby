package object

import (
	"context"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestArrayPush(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := NewArray()
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), array, env)

		result, err := arrayPush(ctx, NewInteger(17))

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewArray(NewInteger(17)), CompareRubyObjectsForTests)
	})

	t.Run("more than one argument", func(t *testing.T) {
		array := NewArray()
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), array, env)

		result, err := arrayPush(ctx, NewInteger(17), NIL, TRUE, FALSE)

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewArray(NewInteger(17), NIL, TRUE, FALSE), CompareRubyObjectsForTests)
	})
}

func TestArrayUnshift(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := NewArray(NewString("first element"))
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), array, env)

		result, err := arrayUnshift(ctx, NewInteger(17))

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewArray(NewInteger(17), NewString("first element")), CompareRubyObjectsForTests)
	})
	t.Run("more than one argument", func(t *testing.T) {
		array := NewArray(NewString("first element"))
		env := env.NewEnvironment[ruby.Object]()
		ctx := call.NewContext[ruby.Object](context.Background(), array, env)

		result, err := arrayUnshift(ctx, NewInteger(17), NIL, TRUE, FALSE)

		assert.NoError(t, err)
		assert.EqualCmpAny(t, result, NewArray(NewInteger(17), NIL, TRUE, FALSE, NewString("first element")), CompareRubyObjectsForTests)
	})
}
