package object

import (
	"reflect"
	"testing"
)

func TestRangeEvalToArray(t *testing.T) {
	t.Run("empty range", func(t *testing.T) {
		rng := &Range{
			Left:      NewInteger(1),
			Right:     NewInteger(1),
			Inclusive: false,
		}
		arr := rng.ToArray()
		expected := &Array{
			Elements: []RubyObject{},
		}
		if !reflect.DeepEqual(arr, expected) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expected, arr)
			t.Fail()
		}
	})
	t.Run("inclusive range", func(t *testing.T) {
		rng := &Range{
			Left:      NewInteger(1),
			Right:     NewInteger(3),
			Inclusive: true,
		}
		arr := rng.ToArray()
		expected := &Array{
			Elements: []RubyObject{
				NewInteger(1),
				NewInteger(2),
				NewInteger(3),
			},
		}
		if !reflect.DeepEqual(arr, expected) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expected, arr)
			t.Fail()
		}
	})
	t.Run("exclusive range", func(t *testing.T) {
		rng := &Range{
			Left:      NewInteger(1),
			Right:     NewInteger(3),
			Inclusive: false,
		}
		arr := rng.ToArray()
		expected := &Array{
			Elements: []RubyObject{
				NewInteger(1),
				NewInteger(2),
			},
		}
		if !reflect.DeepEqual(arr, expected) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expected, arr)
			t.Fail()
		}
	})
}
