package object

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
)

func TestRangeEvalToArray(t *testing.T) {
	t.Run("empty range", func(t *testing.T) {
		rng := &Range{
			Left:      1,
			Right:     1,
			Inclusive: false,
		}
		arr := rng.ToArray()
		utils.AssertEqualCmpAny(t, arr, NewArray(), CompareRubyObjectsForTests)
	})
	t.Run("inclusive range", func(t *testing.T) {
		rng := &Range{
			Left:      1,
			Right:     3,
			Inclusive: true,
		}
		arr := rng.ToArray()
		expected := NewArray(
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		)
		utils.AssertEqualCmpAny(t, arr, expected, CompareRubyObjectsForTests)
	})
	t.Run("exclusive range", func(t *testing.T) {
		rng := &Range{
			Left:      1,
			Right:     3,
			Inclusive: false,
		}
		arr := rng.ToArray()
		expected := NewArray(
			NewInteger(1),
			NewInteger(2),
		)
		utils.AssertEqualCmpAny(t, arr, expected, CompareRubyObjectsForTests)
	})
}
