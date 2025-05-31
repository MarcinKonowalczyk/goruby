package object

import "github.com/MarcinKonowalczyk/goruby/object/ruby"

func swapOrFalse(left, right ruby.Object, swapped bool) bool {
	if swapped {
		// we've already swapped. just return false
		return false
	} else {
		// 1-depth recursive call with swapped arguments
		return rubyObjectsEqual(right, left, true)
	}
}

// TODO: Unify this with rubyObjectsEqual
func CompareRubyObjectsForTests(a, b any) bool {
	// check nils
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// check types
	a_obj, a_ok := a.(ruby.Object)
	// if !ok {
	// 	panic("a is not ruby.Object")
	// }
	b_obj, b_ok := b.(ruby.Object)
	// if !ok {
	// 	panic("b is not ruby.Object")
	// }
	if !a_ok || !b_ok {
		// maybe we're both arrays of ruby.Objects?
		a_arr, a_ok := a.([]ruby.Object)
		b_arr, b_ok := b.([]ruby.Object)
		if a_ok && b_ok {
			// compare the arrays element by element
			if len(a_arr) != len(b_arr) {
				return false
			}
			for i := range a_arr {
				if !CompareRubyObjectsForTests(a_arr[i], b_arr[i]) {
					return false
				}
			}
			return true
		} else {
			if !a_ok {
				panic("a is not RubyObject or []RubyObject")
			}
			if !b_ok {
				panic("b is not RubyObject or []RubyObject")
			}
			panic("a and b are not RubyObject or []RubyObject")
		}
	}

	if a_obj.Class() != b_obj.Class() {
		return false
	}
	// TODO: look into more
	return a_obj.HashKey() == b_obj.HashKey()
	// if a, a_hashable := a_obj.(hashable); a_hashable {
	// 	if b, b_hashable := b_obj.(hashable); b_hashable {
	// 	} else {
	// 		// b is not hashable, we are not equal
	// 		return false
	// 	}
	// }
	// if _, b_hashable := b_obj.(hashable); b_hashable {
	// 	// a is not hashable, we are not equal
	// 	return false
	// }
	// ok, we are not hashable but we are the same class
	// check the addresses

	// addrB := fmt.Sprintf("%p", b_obj)
	// if addrA == addrB {
	// 	return true
	// }
	// return reflect.DeepEqual(a_obj, b_obj)
}

func rubyObjectsEqual(left, right ruby.Object, swapped bool) bool {
	// leftClass := left.Class()
	// rightClass := right.Class()
	// if leftClass != rightClass {
	// 	return swapOrFalse(left, right, swapped)
	// }
	// if left == nil {
	// 	return right == nil || right.Class().Name() == "NilClass"
	// }
	// if left.Class().Name() == "NilClass" {
	// 	return right == nil || right.Class().Name() == "NilClass"
	// }
	// fmt.Printf("left: %T right: %T\n", left, right)
	// fmt.Println("left:", left.Class().Name(), "right:", right.Class().Name())
	switch left := left.(type) {
	case *Integer:
		right_t, ok := safeObjectToInteger(right)
		if !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *Float:
		right_t, ok := safeObjectToFloat(right)
		if !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *String:
		if right_t, ok := right.(*String); !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	case *Array:
		if right_t, ok := right.(*Array); !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			if len(left.Elements) != len(right_t.Elements) {
				return false
			}
			for i, elem := range left.Elements {
				if !rubyObjectsEqual(elem, right_t.Elements[i], swapped) {
					return false
				}
			}
			return true
		}
	case *Hash:
		if right_t, ok := right.(*Hash); !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			if len(left.Map) != len(right_t.Map) {
				return false
			}
			for key, leftValue := range left.ObjectMap() {
				rightValue, ok := right_t.Get(key)
				if !ok {
					return false
				}
				if !rubyObjectsEqual(leftValue, rightValue, swapped) {
					return false
				}
			}

			return true
		}
	case *Symbol:
		if right_t, ok := right.(*Symbol); !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	default:
		return false
	}
}

func RubyObjectsEqual(left, right ruby.Object) bool {
	return rubyObjectsEqual(left, right, false)
}
