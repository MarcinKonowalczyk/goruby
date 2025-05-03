package object

import (
	"fmt"
	"hash/fnv"
)

var (
	symbolClass RubyClassObject = newClass(
		"Symbol",
		symbolMethods,
		symbolClassMethods,
		func(RubyClassObject, ...RubyObject) (RubyObject, error) {
			return &Symbol{}, nil
		},
	)
	TRUE  RubyObject = &Symbol{Value: "true"}
	FALSE RubyObject = &Symbol{Value: "false"}
	NIL   RubyObject = &Symbol{Value: "nil"}
)

func init() {
	CLASSES.Set("Symbol", symbolClass)
}

type Symbol struct {
	Value string
}

// Inspect returns the value of the symbol
func (s *Symbol) Inspect() string { return ":" + s.Value }

// Type returns SYMBOL_OBJ
func (s *Symbol) Type() Type { return SYMBOL_OBJ }

// Class returns symbolClass
func (s *Symbol) Class() RubyClass { return symbolClass }

func (s *Symbol) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hashKey{Type: s.Type(), Value: h.Sum64()}
}

var symbolClassMethods = map[string]RubyMethod{}

var symbolMethods = map[string]RubyMethod{
	"to_s": withArity(0, newMethod(symbolToS)),
	"to_i": withArity(0, newMethod(symbolToI)),
	"size": withArity(0, newMethod(symbolSize)),
}

func symbolToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	if sym, ok := context.Receiver().(*Symbol); ok {
		return &String{Value: sym.Value}, nil
	}
	return nil, nil
}

func symbolSize(context CallContext, args ...RubyObject) (RubyObject, error) {
	if sym, ok := context.Receiver().(*Symbol); ok {
		return NewInteger(int64(len(sym.Value))), nil
	}
	return nil, nil
}

func symbolToI(context CallContext, args ...RubyObject) (RubyObject, error) {
	if boolean, ok := SymbolToBool(context.Receiver()); ok {
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
func SymbolToBool(o RubyObject) (val bool, ok bool) {
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
			fmt.Println("SymbolToBool: unknown symbol value:", sym.Value)
			return false, false
		}
	}
	return false, false
}

// If a symbol is "nil", returns ok as 'true'.
func SymbolToNil(o RubyObject) (ok bool) {
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
