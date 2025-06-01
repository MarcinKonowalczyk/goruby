package interpreter_test

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/pipelines/interpreter"
	"github.com/MarcinKonowalczyk/goruby/testutils/assert"
)

func TestInterpreterInterpret(t *testing.T) {
	t.Run("return proper result", func(t *testing.T) {
		input := `
			def foo
				3
			end

			x = 5

			def add x, y
				x + y
			end

			add foo, x
			`
		i := interpreter.NewBasicInterpreter()

		out, err := i.InterpretCode(input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		assert.That(t, ok, "Expected *object.Integer, got %T\n", out)
		assert.Equal(t, res.Value, 8)
	})
}
