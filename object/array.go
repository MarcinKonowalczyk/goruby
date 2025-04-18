package object

import (
	"hash/fnv"
	"strings"
)

var arrayClass RubyClassObject = newClass(
	"Array",
	objectClass,
	arrayMethods,
	arrayClassMethods,
	func(c RubyClassObject, args ...RubyObject) (RubyObject, error) { return NewArray(args...), nil },
)

func init() {
	classes.Set("Array", arrayClass)
}

// NewArray returns a new array populated with elements.
func NewArray(elements ...RubyObject) *Array {
	arr := &Array{Elements: make([]RubyObject, len(elements))}
	for i, elem := range elements {
		arr.Elements[i] = elem
	}
	return arr
}

// An Array represents a Ruby Array
type Array struct {
	Elements []RubyObject
}

// Type returns the ObjectType of the array
func (a *Array) Type() Type { return ARRAY_OBJ }

// Inspect returns all elements within the array, divided by comma and
// surrounded by brackets
func (a *Array) Inspect() string {
	elems := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		elems[i] = elem.Inspect()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

// Class returns the class of the Array
func (a *Array) Class() RubyClass { return arrayClass }
func (a *Array) hashKey() hashKey {
	h := fnv.New64a()
	for _, e := range a.Elements {
		h.Write(hash(e).bytes())
	}
	return hashKey{Type: a.Type(), Value: h.Sum64()}
}

var arrayClassMethods = map[string]RubyMethod{}

var arrayMethods = map[string]RubyMethod{
	"push":     publicMethod(arrayPush),
	"unshift":  publicMethod(arrayUnshift),
	"size":     publicMethod(arraySize),
	"length":   publicMethod(arraySize),
	"find_all": publicMethod(arrayFindAll),
	"first":    publicMethod(arrayFirst),
	"map":      publicMethod(arrayMap),
	"all?":     publicMethod(arrayAll),
	"join":     publicMethod(arrayJoin),
	"include?": publicMethod(arrayInclude),
	"each":     publicMethod(arrayEach),
	"reject":   publicMethod(arrayReject),
	"pop":      publicMethod(arrayPop),
}

func arrayPush(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	array.Elements = append(array.Elements, args...)
	return array, nil
}

func arrayUnshift(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	array.Elements = append(args, array.Elements...)
	return array, nil
}

func arraySize(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	return &Integer{Value: int64(len(array.Elements))}, nil
}

func arrayFindAll(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("find_all requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("find_all requires a block")
	}
	result := NewArray()
	for _, elem := range array.Elements {
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

func arrayFirst(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		if len(array.Elements) == 0 {
			return nil, NewArgumentError("array is empty")
		}
		return array.Elements[0], nil
	}
	if len(args) > 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 0..1)", len(args))
	}
	n, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError("argument must be an Integer")
	}
	if n.Value < 0 {
		return nil, NewArgumentError("negative array size (or size too big)")
	}
	if int(n.Value) > len(array.Elements) {
		return nil, NewArgumentError("array size too big")
	}
	result := NewArray()
	for i := 0; i < int(n.Value); i++ {
		result.Elements = append(result.Elements, array.Elements[i])
	}
	return result, nil
}

func arrayMap(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	result := NewArray()
	for _, elem := range array.Elements {
		ret, err := proc.Call(context, elem)
		if err != nil {
			return nil, err
		}
		result.Elements = append(result.Elements, ret)
	}
	return result, nil
}

func arrayAll(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("all? requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("all? requires a block")
	}
	for _, elem := range array.Elements {
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

func arrayJoin(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("join requires at least 1 argument")
	}
	separator, ok := args[0].(*String)
	if !ok {
		return nil, NewArgumentError("argument must be a String")
	}
	element_strings := make([]string, len(array.Elements))
	for i, elem := range array.Elements {
		element_strings[i] = elem.Inspect()
	}
	result := strings.Join(element_strings, separator.Value)
	return &String{Value: result}, nil
}

func arrayInclude(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("include? requires at least 1 argument")
	}
	// compare the first argument with all elements in the array
	arg := args[0]
	for _, elem := range array.Elements {
		// if elem.Class() == arg.Class() {
		// 	if elem.Inspect() == arg.Inspect() {
		// 		return TRUE, nil
		// 	}
		// 	return FALSE, nil
		// } else {
		// 	if _, ok := elem.(*Integer); ok {
		// 		// also compare as if it was a float

		// 	}
		// }
		ctx := NewCallContext(context.Env(), elem)
		ret, err := Send(ctx, "==", arg)
		if err != nil {
			// fmt.Println("Error in arrayInclude:", err)
			continue
			// return nil, err
		}
		boolean := ret.(*Boolean)
		if boolean.Value {
			return TRUE, nil
		}
	}
	return FALSE, nil
}

func arrayEach(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	for _, elem := range array.Elements {
		_, err := proc.Call(context, elem)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func arrayReject(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Proc)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	result := NewArray()
	for _, elem := range array.Elements {
		ret, err := proc.Call(context, elem)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, NewArgumentError("map requires a block to return a boolean")
		}
		if ret.Type() != BOOLEAN_OBJ {
			return nil, NewArgumentError("map requires a block to return a boolean")
		}
		boolean, ok := ret.(*Boolean)
		if !ok {
			return nil, NewArgumentError("map requires a block to return a boolean")
		}
		if !boolean.Value {
			result.Elements = append(result.Elements, elem)
		}
	}
	return result, nil
}

func arrayPop(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	if len(array.Elements) == 0 {
		return NIL, nil
	}
	elem := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]
	return elem, nil
}
