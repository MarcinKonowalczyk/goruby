package object_test

import (
	"context"
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestWithArity(t *testing.T) {
	wrappedMethod := ruby.NewMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
		return object.NewInteger(1), nil
	})

	tests := []struct {
		arity     int
		arguments []ruby.Object
		result    ruby.Object
		err       error
	}{
		{
			1,
			[]ruby.Object{object.NIL},
			object.NewInteger(1),
			nil,
		},
		{
			1,
			[]ruby.Object{object.NIL, object.NIL},
			nil,
			object.NewWrongNumberOfArgumentsError(1, 2),
		},
		{
			2,
			[]ruby.Object{object.NIL},
			nil,
			object.NewWrongNumberOfArgumentsError(2, 1),
		},
	}

	for _, testCase := range tests {
		fn := object.WithArity(testCase.arity, wrappedMethod)
		ctx := call.NewContext[ruby.Object](context.Background()).WithReceiver(object.NIL)
		result, err := fn.Call(ctx, testCase.arguments...)
		assert.EqualCmpAny(t, result, testCase.result, object.CompareRubyObjectsForTests)
		assert.Error(t, err, testCase.err)
	}
}
