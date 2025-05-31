package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
	"github.com/pkg/errors"
)

type testRubyObject struct {
	class ruby.ClassObject
	Name  string
}

func (t *testRubyObject) Inspect() string { return "TEST OBJECT" }
func (t *testRubyObject) Class() ruby.Class {
	if t.class != nil {
		return t.class
	}
	return bottomClass
}
func (t *testRubyObject) HashKey() hash.Key {
	return hash.Key(99)
}

func TestSend(t *testing.T) {
	// }
	methods := map[string]ruby.Method{
		"a_method": ruby.NewMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return TRUE, nil
		}),
		"another_method": ruby.NewMethod(func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			return FALSE, nil
		}),
	}
	t.Run("normal object as context", func(t *testing.T) {
		// ctx := &callContext{
		// 	receiver: &testRubyObject{
		// 		class: &class{
		// 			name:            "base class",
		// 			instanceMethods: NewMethodSet(methods),
		// 		},
		// 	},
		// }
		ctx := call.NewContext[ruby.Object](
			&testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: ruby.NewMethodSet(methods),
				},
			},
			nil,
		)

		tests := []struct {
			method         string
			expectedResult ruby.Object
			expectedError  error
		}{
			{
				"a_method",
				TRUE,
				nil,
			},
			{
				"another_method",
				FALSE,
				nil,
			},
			{
				"unknown_method",
				nil,
				NewNoMethodError(ctx.Receiver(), "unknown_method"),
			},
		}

		for _, testCase := range tests {
			result, err := Send(ctx, testCase.method)

			assert.Error(t, errors.Cause(err), testCase.expectedError)
			assert.EqualCmpAny(t, result, testCase.expectedResult, CompareRubyObjectsForTests)
		}
	})
}
