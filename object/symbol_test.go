package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestSymbol_hashKey(t *testing.T) {
	hello1 := NewSymbol("Hello World")
	hello2 := NewSymbol("Hello World")
	diff1 := NewSymbol("My name is johnny")
	diff2 := NewSymbol("My name is johnny")

	utils.AssertEqual(t, hello1.HashKey(), hello2.HashKey())
	utils.AssertEqual(t, diff1.HashKey(), diff2.HashKey())
	utils.AssertNotEqual(t, hello1.HashKey(), diff1.HashKey())
}

func TestSymbolToS(t *testing.T) {
	context := &callContext{
		receiver: NewSymbol("foo"),
	}

	result, err := symbolToS(context, nil)

	utils.AssertNoError(t, err)

	expected := NewString("foo")

	utils.AssertEqualCmpAny(t, result, expected, CompareRubyObjectsForTests)
}

func TestSymbolToBool(t *testing.T) {
	t.Run("true object", func(t *testing.T) {
		val, ok := SymbolToBool(TRUE)
		utils.Assert(t, ok)
		utils.Assert(t, val)
	})

	t.Run("false object", func(t *testing.T) {
		val, ok := SymbolToBool(FALSE)
		utils.Assert(t, ok)
		utils.Assert(t, !val)
	})

	t.Run("some other symbol", func(t *testing.T) {
		val, ok := SymbolToBool(NewSymbol("foo"))
		utils.Assert(t, !ok)
		utils.Assert(t, !val)
	})

	t.Run("some non-symbol object", func(t *testing.T) {
		val, ok := SymbolToBool(NewString("foo"))
		utils.Assert(t, !ok)
		utils.Assert(t, !val)
	})

	t.Run("nil", func(t *testing.T) {
		val, ok := SymbolToBool(nil)
		utils.Assert(t, !ok)
		utils.Assert(t, !val)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var b *Symbol = nil
		val, ok := SymbolToBool(b)
		utils.Assert(t, !ok)
		utils.Assert(t, !val)
	})
}
