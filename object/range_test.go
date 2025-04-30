package object

import (
	"reflect"
	"testing"
)

func TestRangeEvalToArray(t *testing.T) {
	t.Run("empty range", func(t *testing.T) {
		rng := &Range{
			Left:      &Integer{Value: 1},
			Right:     &Integer{Value: 1},
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
			Left:      &Integer{Value: 1},
			Right:     &Integer{Value: 3},
			Inclusive: true,
		}
		arr := rng.ToArray()
		expected := &Array{
			Elements: []RubyObject{
				&Integer{Value: 1},
				&Integer{Value: 2},
				&Integer{Value: 3},
			},
		}
		if !reflect.DeepEqual(arr, expected) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expected, arr)
			t.Fail()
		}
	})
	t.Run("exclusive range", func(t *testing.T) {
		rng := &Range{
			Left:      &Integer{Value: 1},
			Right:     &Integer{Value: 3},
			Inclusive: false,
		}
		arr := rng.ToArray()
		expected := &Array{
			Elements: []RubyObject{
				&Integer{Value: 1},
				&Integer{Value: 2},
			},
		}
		if !reflect.DeepEqual(arr, expected) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expected, arr)
			t.Fail()
		}
	})
}
