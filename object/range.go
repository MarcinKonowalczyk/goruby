package object

import (
	"hash/fnv"
	"strings"
)

var rangeClass RubyClassObject = newClass(
	"Range",
	rangeMethods,
	nil,
	notInstantiatable,
)

func init() {
	CLASSES.Set("Range", rangeClass)
}

type Range struct {
	Left      *Integer
	Right     *Integer
	Inclusive bool
}

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

func (a *Range) Class() RubyClass { return rangeClass }

func (a *Range) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(a.Left.Inspect()))
	h.Write([]byte(a.Right.Inspect()))
	if a.Inclusive {
		h.Write([]byte("1"))
	} else {
		h.Write([]byte("0"))
	}
	return HashKey(h.Sum64())
}

var rangeMethods = map[string]RubyMethod{
	"find_all": newMethod(rangeFindAll),
	"all?":     newMethod(rangeAll),
}

func (rang *Range) ToArray() *Array {
	result := NewArray()
	left, right := rang.Left.Value, rang.Right.Value
	if rang.Inclusive {
		right++
	}
	for i := left; i < right; i++ {
		result.Elements = append(result.Elements, NewInteger(i))
	}
	return result
}

func rangeFindAll(context CallContext, args ...RubyObject) (RubyObject, error) {
	rng, _ := context.Receiver().(*Range)
	if len(args) == 0 {
		return nil, NewArgumentError("find_all requires a block")
	}
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("find_all requires a block")
	}
	self, _ := context.Env().Get("self")
	self_class := self.Class()
	fn, ok := self_class.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(self, proc.Value)
	}
	// evaluate the range
	result := NewArray()
	for _, elem := range rng.ToArray().Elements {
		ret, err := fn.Call(context, elem)
		if err != nil {
			return nil, err
		}
		if ret == nil {
			return nil, NewArgumentError("find_all requires a block to return a boolean")
		}
		val, ok := SymbolToBool(ret)
		if !ok {
			return nil, NewArgumentError("find_all requires a block to return a boolean")
		}
		if val {
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
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("all? requires a block")
	}
	self, _ := context.Env().Get("self")
	self_class := self.Class()
	fn, ok := self_class.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(self, proc.Value)
	}
	for _, elem := range rng.ToArray().Elements {
		ret, err := fn.Call(context, elem)
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
