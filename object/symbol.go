package object

import (
	"fmt"
	"hash/fnv"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var (
	symbolClass ruby.ClassObject = nil_class

	// unique symbols
	TRUE  ruby.Object = &Symbol{Value: "true"}
	FALSE ruby.Object = &Symbol{Value: "false"}
	NIL   ruby.Object = &Symbol{Value: "nil"}
	FUNCS ruby.Object = &Symbol{Value: "funcs"}
)

// instantiate the symbol class
//
//go:inline
func initSymbolClass() {
	symbolClass = newClass("Symbol", symbolMethods, symbolClassMethods)
}

func init() {
	initSymbolClass()
	CLASSES.Set("Symbol", symbolClass)
}

//go:inline
func NewSymbol(value string) *Symbol {
	switch value {
	case "true":
		return TRUE.(*Symbol)
	case "false":
		return FALSE.(*Symbol)
	case "nil":
		return NIL.(*Symbol)
	case "funcs":
		return FUNCS.(*Symbol)
	default:
		return &Symbol{Value: value}
	}
}

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string { return ":" + s.Value }

func (s *Symbol) Class() ruby.Class {
	if symbolClass == nil {
		panic("symbolClass is nil")
	}
	if sc, ok := symbolClass.(*class); ok && sc == nil_class {
		initSymbolClass()
		if sc2, ok := symbolClass.(*class); ok && sc2 == nil_class {
			panic("symbolClass is *still* nil_class")
		}
	}
	return symbolClass
}

func (s *Symbol) HashKey() hash.Key {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hash.Key(h.Sum64())
}

var (
	_ ruby.Object = &Symbol{}
	_ ruby.Class  = symbolClass
)

var symbolClassMethods = map[string]ruby.Method{}

var symbolMethods = map[string]ruby.Method{
	"to_s": WithArity(0, ruby.NewMethod(symbolToS)),
	"to_i": WithArity(0, ruby.NewMethod(symbolToI)),
	"size": WithArity(0, ruby.NewMethod(symbolSize)),
}

func symbolToS(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	if sym, ok := ctx.Receiver().(*Symbol); ok {
		return NewString(sym.Value), nil
	}
	return nil, nil
}

func symbolSize(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	if sym, ok := ctx.Receiver().(*Symbol); ok {
		return NewInteger(int64(len(sym.Value))), nil
	}
	return nil, nil
}

func symbolToI(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	if boolean, ok := SymbolToBool(ctx.Receiver()); ok {
		if boolean {
			return NewInteger(1), nil
		}
		return NewInteger(0), nil
	} else {
		return nil, NewTypeError("Symbol to_i: expected boolean")
	}
}

// If a symbol is "true" or "false", it returns the corresponding boolean value.
// Otherwise, it returns false and ok as false.
func SymbolToBool(o ruby.Object) (val bool, ok bool) {
	if sym, ok := o.(*Symbol); ok {
		if sym == nil {
			// nil pointer, not ok
			return false, false
		}
		if sym.Value == "true" {
			return true, true
		} else if sym.Value == "false" {
			return false, true
		} else {
			return false, false
		}
	}
	return false, false
}

// If a symbol is "nil", returns ok as 'true'.
func SymbolToNil(o ruby.Object) (ok bool) {
	if sym, ok := o.(*Symbol); ok {
		if sym == nil {
			// nil pointer, not ok
			// *must be a ruby object!*
			return false
		}
		if sym.Value == "nil" {
			return true
		} else {
			fmt.Println("SymbolToNil: unknown symbol value:", sym.Value)
			return false
		}
	}
	return false
}
