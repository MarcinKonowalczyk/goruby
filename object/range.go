package object

import (
	"hash/fnv"
	"strings"
)

var rangeClass RubyClassObject = newClass(
	"Range",
	objectClass,
	rangeMethods,
	rangeClassMethods,
	notInstantiatable,
)

func init() {
	classes.Set("Range", rangeClass)
}

// An Range represents a Ruby Range
type Range struct {
	Left      *Integer
	Right     *Integer
	Inclusive bool
}

// Type returns the ObjectType of the range
func (a *Range) Type() Type { return RANGE_OBJ }

// Inspect returns all elements within the range, divided by comma and
// surrounded by brackets
func (a *Range) Inspect() string {
	var out strings.Builder
	out.WriteString(a.Left.Inspect())
	if a.Inclusive {
		out.WriteString("..")
	} else {
		out.WriteString("...")
	}
	out.WriteString(a.Right.Inspect())
	return out.String()
}

// Class returns the class of the Range
func (a *Range) Class() RubyClass { return rangeClass }
func (a *Range) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(a.Left.Inspect()))
	h.Write([]byte(a.Right.Inspect()))
	if a.Inclusive {
		h.Write([]byte("1"))
	} else {
		h.Write([]byte("0"))
	}
	return hashKey{Type: a.Type(), Value: h.Sum64()}
}

var rangeClassMethods = map[string]RubyMethod{}

var rangeMethods = map[string]RubyMethod{
	// "push":     publicMethod(rangePush),
	// "unshift":  publicMethod(rangeUnshift),
	// "size":     publicMethod(rangeSize),
	"find_all": publicMethod(rangeFindAll),
	// "first":    publicMethod(rangeFirst),
	// "map":      publicMethod(rangeMap),
	"all?": publicMethod(rangeAll),
	// "join":     publicMethod(rangeJoin),
}

// func rangePush(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	range, _ := context.Receiver().(*Range)
// 	range.Elements = append(range.Elements, args...)
// 	return range, nil
// }

// func rangeUnshift(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	range, _ := context.Receiver().(*Range)
// 	range.Elements = append(args, range.Elements...)
// 	return range, nil
// }

// func rangeSize(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	range, _ := context.Receiver().(*Range)
// 	return &Integer{Value: int64(len(range.Elements))}, nil
// }

func (rang *Range) ToArray() *Array {
	result := NewArray()
	left, right := rang.Left.Value, rang.Right.Value
	if rang.Inclusive {
		right++
	}
	for i := left; i < right; i++ {
		result.Elements = append(result.Elements, &Integer{Value: i})
	}
	return result
}

func rangeFindAll(context CallContext, args ...RubyObject) (RubyObject, error) {
	rng, _ := context.Receiver().(*Range)
	if len(args) == 0 {
		return nil, NewArgumentError("find_all requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("find_all requires a block")
	}
	// evaluate the range
	result := NewArray()
	for _, elem := range rng.ToArray().Elements {
		ret, err := proc.Call(context, elem)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, NewArgumentError("find_all requires a block to return a boolean")
		}
		if ret.Type() != BOOLEAN_OBJ {
			return nil, NewArgumentError("find_all requires a block to return a boolean")
		}
		boolean, ok := ret.(*Boolean)
		if !ok {
			return nil, NewArgumentError("find_all requires a block to return a boolean")
		}
		if boolean.Value {
			result.Elements = append(result.Elements, elem)
		}
	}
	return result, nil
}

func rangeAll(context CallContext, args ...RubyObject) (RubyObject, error) {
	rng, _ := context.Receiver().(*Range)
	if len(args) == 0 {
		return nil, NewArgumentError("all? requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("all? requires a block")
	}
	for _, elem := range rng.ToArray().Elements {
		ret, err := proc.Call(context, elem)
		if err != nil {
			return nil, err
		}
		if ret.Type() != BOOLEAN_OBJ {
			return nil, NewArgumentError("all? requires a block to return a boolean")
		}
		boolean, ok := ret.(*Boolean)
		if !ok {
			return nil, NewArgumentError("all? requires a block to return a boolean")
		}
		if !boolean.Value {
			return FALSE, nil
		}
	}
	return TRUE, nil
}

// func rangeJoin(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	range, _ := context.Receiver().(*Range)
// 	if len(args) == 0 {
// 		return nil, NewArgumentError("join requires at least 1 argument")
// 	}
// 	separator, ok := args[0].(*String)
// 	if !ok {
// 		return nil, NewArgumentError("argument must be a String")
// 	}
// 	element_strings := make([]string, len(range.Elements))
// 	for i, elem := range range.Elements {
// 		element_strings[i] = elem.Inspect()
// 	}
// 	result := strings.Join(element_strings, separator.Value)
// 	return &String{Value: result}, nil
// }
