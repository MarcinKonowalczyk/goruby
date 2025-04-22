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

// func TestArrayUnshift(t *testing.T) {
// 	t.Run("one argument", func(t *testing.T) {
// 		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
// 		env := NewEnvironment()
// 		context := &callContext{
// 			receiver: array,
// 			env:      env,
// 		}

// 		result, err := arrayUnshift(context, &Integer{Value: 17})

// 		checkError(t, err, nil)

// 		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}, &String{Value: "first element"}}}

// 		checkResult(t, result, expectedResult)

// 		if !reflect.DeepEqual(expectedResult, array) {
// 			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
// 			t.Fail()
// 		}
// 	})
// 	t.Run("more than one argument", func(t *testing.T) {
// 		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
// 		env := NewEnvironment()
// 		context := &callContext{
// 			receiver: array,
// 			env:      env,
// 		}

// 		result, err := arrayUnshift(context, &Integer{Value: 17}, NIL, TRUE, FALSE)

// 		checkError(t, err, nil)

// 		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}, NIL, TRUE, FALSE, &String{Value: "first element"}}}

// 		checkResult(t, result, expectedResult)

// 		if !reflect.DeepEqual(expectedResult, array) {
// 			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
// 			t.Fail()
// 		}
// 	})
// }
