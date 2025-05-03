package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestArrayPush(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := &Array{}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayPush(context, &Integer{Value: 17})

		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewArray(NewInteger(17)), CompareRubyObjectsForTests)
	})

	t.Run("more than one argument", func(t *testing.T) {
		array := &Array{}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayPush(context, &Integer{Value: 17}, NIL, TRUE, FALSE)

		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewArray(NewInteger(17), NIL, TRUE, FALSE), CompareRubyObjectsForTests)
	})
}

func TestArrayUnshift(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayUnshift(context, &Integer{Value: 17})

		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewArray(NewInteger(17), NewString("first element")), CompareRubyObjectsForTests)
	})
	t.Run("more than one argument", func(t *testing.T) {
		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayUnshift(context, &Integer{Value: 17}, NIL, TRUE, FALSE)

		utils.AssertNoError(t, err)
		utils.AssertEqualCmpAny(t, result, NewArray(NewInteger(17), NIL, TRUE, FALSE, NewString("first element")), CompareRubyObjectsForTests)
	})
}
