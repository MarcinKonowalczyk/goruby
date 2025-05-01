package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetMethods(t *testing.T) {
	superClassMethods := map[string]RubyMethod{
		"super_foo": publicMethod(nil),
		"super_bar": publicMethod(nil),
	}
	contextMethods := map[string]RubyMethod{
		"foo": publicMethod(nil),
		"bar": publicMethod(nil),
	}
	classWithoutSuperclass := &class{
		instanceMethods: NewMethodSet(contextMethods),
		superClass:      nil,
	}
	classWithSuperclass := &class{
		instanceMethods: NewMethodSet(contextMethods),
		superClass: &class{
			instanceMethods: NewMethodSet(superClassMethods),
			superClass:      nil,
		},
	}

	tests := []struct {
		name                 string
		class                RubyClass
		addSuperclassMethods bool
		expectedMethods      []string
	}{
		{
			"no superclass methods add super methods",
			classWithoutSuperclass,
			true,
			[]string{":foo", ":bar"},
		},
		{
			"no superclass methods add no super methods",
			classWithoutSuperclass,
			false,
			[]string{":foo", ":bar"},
		},
		{
			"with superclass methods add super methods",
			classWithSuperclass,
			true,
			[]string{":foo", ":bar", ":super_foo", ":super_bar"},
		},
		{
			"with superclass methods dont add super methods",
			classWithSuperclass,
			false,
			[]string{":foo", ":bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMethods(tt.class, tt.addSuperclassMethods)

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

			expectedLen := len(tt.expectedMethods)

			if len(result.Elements) != expectedLen {
				t.Logf("Expected %d items, got %d", expectedLen, len(result.Elements))
				t.Fail()
			}

			sort.Strings(tt.expectedMethods)
			sort.Strings(methods)

			if !reflect.DeepEqual(tt.expectedMethods, methods) {
				t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", tt.expectedMethods, methods)
				t.Fail()
			}
		})
	}
}
