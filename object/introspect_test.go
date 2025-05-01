package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetMethods(t *testing.T) {
	contextMethods := map[string]RubyMethod{
		"foo": publicMethod(nil),
		"bar": publicMethod(nil),
	}
	myClass := &class{
		instanceMethods: NewMethodSet(contextMethods),
	}

	t.Run("no super methods", func(t *testing.T) {
		result := getMethods(myClass, false)
		var methods []string
		for i, elem := range result.Elements {
			sym, ok := elem.(*Symbol)
			if !ok {
				t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
				t.Fail()
			} else {
				methods = append(methods, sym.Inspect())
			}
		}
		expectedLen := len(contextMethods)
		if len(result.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(result.Elements))
			t.Fail()
		}
		sort.Strings(methods)
		expectedMethods := []string{":foo", ":bar"}
		sort.Strings(expectedMethods)
		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})

	t.Run("add super methods", func(t *testing.T) {
		result := getMethods(myClass, true)
		var methods []string
		for i, elem := range result.Elements {
			sym, ok := elem.(*Symbol)
			if !ok {
				t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
				t.Fail()
			} else {
				methods = append(methods, sym.Inspect())
			}
		}
		someExpectedMethods := []string{":foo", ":bar", ":raise", ":==", ":!="}
		sort.Strings(someExpectedMethods)
		for _, method := range someExpectedMethods {
			var found bool = false
			for _, m := range methods {
				if method == m {
					found = true
					break
				}
			}
			if !found {
				t.Logf("Expected method %s to be found in %s", method, methods)
				t.Fail()
			}
		}
	})
}
