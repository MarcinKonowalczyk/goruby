package interpreter_test

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/interpreter"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/utils"
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
		i := interpreter.NewInterpreter()

		out, err := i.Interpret("", input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		utils.Assert(t, ok, "Expected *object.Integer, got %T\n", out)
		utils.AssertEqual(t, res.Value, 8)
	})
}
