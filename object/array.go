package object

import (
	"hash/fnv"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var arrayClass ruby.ClassObject = newClass(
	"Array",
	arrayMethods,
	nil,
)

func init() {
	CLASSES.Set("Array", arrayClass)
}

// NewArray returns a new array populated with elements.
func NewArray(elements ...ruby.Object) *Array {
	arr := &Array{Elements: make([]ruby.Object, len(elements))}
	if len(elements) == 0 {
		return arr
	}
	copy(arr.Elements, elements)
	return arr
}

type Array struct {
	Elements []ruby.Object
}

func (a *Array) Inspect() string {
	elems := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		var elem_str string
		switch elem.(type) {
		case *String:
			elem_str = "\"" + elem.Inspect() + "\""
		default:
			elem_str = elem.Inspect()
		}
		elems[i] = elem_str
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

func (a *Array) Class() ruby.Class { return arrayClass }
func (a *Array) HashKey() hash.Key {
	h := fnv.New64a()
	for _, e := range a.Elements {
		h.Write(e.HashKey().Bytes())
	}
	return hash.Key(h.Sum64())
}

var arrayMethods = map[string]ruby.Method{
	"push":     ruby.NewMethod(arrayPush),
	"unshift":  ruby.NewMethod(arrayUnshift),
	"size":     ruby.NewMethod(arraySize),
	"length":   ruby.NewMethod(arraySize),
	"find_all": ruby.NewMethod(arrayFindAll),
	"first":    ruby.NewMethod(arrayFirst),
	"map":      ruby.NewMethod(arrayMap),
	"all?":     ruby.NewMethod(arrayAll),
	"join":     ruby.NewMethod(arrayJoin),
	"include?": ruby.NewMethod(arrayInclude),
	"each":     ruby.NewMethod(arrayEach),
	"reject":   ruby.NewMethod(arrayReject),
	"pop":      ruby.NewMethod(arrayPop),
	"-":        ruby.NewMethod(arrayMinus),
	"+":        ruby.NewMethod(arrayPlus),
	"*":        ruby.NewMethod(arrayAst),
}

func arrayPush(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	array.Elements = append(array.Elements, args...)
	return array, nil
}

func arrayUnshift(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	array.Elements = append(args, array.Elements...)
	return array, nil
}

func arraySize(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	return NewInteger(int64(len(array.Elements))), nil
}

func arrayFindAll(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("(1) array find_all requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("(2) array find_all requires a block")
	}
	fn, ok := FUNCS_STORE.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(FUNCS_STORE, proc.Value)
	}
	result := NewArray()
	for _, element := range array.Elements {
		ret, err := fn.Call(ctx, element)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, NewArgumentError("find_all requires a block to return a boolean, not nil")
		}
		val, ok := SymbolToBool(ret)
		if !ok {
			return nil, NewArgumentError("find_all requires a block to return a boolean, not %s", ret.Inspect())
		}
		if val {
			result.Elements = append(result.Elements, element)
		}
	}
	return result, nil
}

func arrayFirst(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
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

func arrayMap(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	trace.MessageCtx(ctx, func() string { return ctx.Receiver().Inspect() })

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	fn, ok := FUNCS_STORE.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(FUNCS_STORE, proc.Value)
	}
	result := NewArray()
	for _, elem := range array.Elements {
		ret, err := fn.Call(ctx, elem)
		if err != nil {
			return nil, err
		}
		result.Elements = append(result.Elements, ret)
	}
	return result, nil
}

func arrayAll(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("all? requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("all? requires a block")
	}
	fn, ok := FUNCS_STORE.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(FUNCS_STORE, proc.Value)
	}
	for _, elem := range array.Elements {
		ret, err := fn.Call(ctx, elem)
		if err != nil {
			return nil, err
		}
		val, ok := SymbolToBool(ret)
		if !ok {
			return nil, NewArgumentError("all? requires a block to return a boolean")
		}
		if !val {
			return FALSE, nil
		}
	}
	return TRUE, nil
}

func arrayJoin(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
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
	return NewString(result), nil
}

func arrayInclude(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
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
		ret, err := Send(ctx.WithReceiver(elem), "==", arg)
		if err != nil {
			// fmt.Println("Error in arrayInclude:", err)
			continue
			// return nil, err
		}
		val, ok := SymbolToBool(ret)
		if !ok {
			return nil, NewArgumentError("include? requires a block to return a boolean")
		}
		if val {
			return TRUE, nil
		}
	}
	return FALSE, nil
}

func arrayEach(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	fn, ok := FUNCS_STORE.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(FUNCS_STORE, proc.Value)
	}
	for _, elem := range array.Elements {
		_, err := fn.Call(ctx, elem)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func arrayReject(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("map requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	fn, ok := FUNCS_STORE.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(FUNCS_STORE, proc.Value)
	}
	result := NewArray()
	for _, elem := range array.Elements {
		ret, err := fn.Call(ctx, elem)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, NewArgumentError("map requires a block to return a boolean")
		}
		val, ok := SymbolToBool(ret)
		if !ok {
			return nil, NewArgumentError("map requires a block to return a boolean")
		}
		if !val {
			result.Elements = append(result.Elements, elem)
		}
	}
	return result, nil
}

func arrayPop(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(array.Elements) == 0 {
		return NIL, nil
	}
	elem := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]
	return elem, nil
}

func arrayMinus(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("array minus requires at least 1 argument")
	}
	otherArray, ok := args[0].(*Array)
	if !ok {
		return nil, NewArgumentError("argument must be an Array")
	}
	result := NewArray()
	for _, elem := range array.Elements {
		include := false
		for _, otherElem := range otherArray.Elements {
			ret, err := Send(ctx.WithReceiver(elem), "==", otherElem)
			if err != nil {
				return nil, err
			}
			val, ok := SymbolToBool(ret)
			if !ok {
				return nil, NewArgumentError("array minus requires a block to return a boolean")
			}
			if val {
				include = true
				break
			}
		}
		if !include {
			result.Elements = append(result.Elements, elem)
		}
	}
	return result, nil
}

func arrayPlus(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("array plus requires at least 1 argument")
	}
	otherArray, ok := args[0].(*Array)
	if !ok {
		return nil, NewArgumentError("argument must be an Array")
	}
	result := NewArray()
	result.Elements = append(result.Elements, array.Elements...)
	result.Elements = append(result.Elements, otherArray.Elements...)
	return result, nil
}

func arrayAst(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()

	array, _ := ctx.Receiver().(*Array)
	if len(args) == 0 {
		return nil, NewArgumentError("array ast requires at least 1 argument")
	}
	switch arg := args[0].(type) {
	case *Integer:
		// repeat the array n times
		if arg.Value < 0 {
			return nil, NewArgumentError("negative array size (or size too big)")
		}
		// repeat the array n times
		result := NewArray()
		for i := int64(0); i < arg.Value; i++ {
			result.Elements = append(result.Elements, array.Elements...)
		}
		// fmt.Println("arrayAst", result.Inspect())
		return result, nil
	case *Float:
		// repeat the array floor(n) times
		if arg.Value < 0 {
			return nil, NewArgumentError("negative array size (or size too big)")
		}
		times := int64(arg.Value)
		result := NewArray()
		for i := int64(0); i < times; i++ {
			result.Elements = append(result.Elements, array.Elements...)
		}
		return result, nil
	case *String:
		// [1, 2,3] * ',' id the same as [1, 2,3].join(',')

		joiner := arg.Value
		element_strings := make([]string, len(array.Elements))
		for i, elem := range array.Elements {
			element_strings[i] = elem.Inspect()
		}
		result := strings.Join(element_strings, joiner)
		return NewString(result), nil

	default:
		return nil, NewArgumentError("argument must be an Integer, or a String")
	}
}
