package object

import (
	"fmt"
	"math"
)

var integerClass RubyClassObject = newClass(
	"Integer", objectClass, integerMethods, integerClassMethods, notInstantiatable,
)

func init() {
	classes.Set("Integer", integerClass)
}

// NewInteger returns a new Integer with the given value
func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

// Integer represents an integer in Ruby
type Integer struct {
	Value int64
}

// Inspect returns the value as string
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns INTEGER_OBJ
func (i *Integer) Type() Type { return INTEGER_OBJ }

// Class returns integerClass
func (i *Integer) Class() RubyClass { return integerClass }

var (
	_ RubyObject = &Integer{}
)

func (i *Integer) hashKey() hashKey {
	return hashKey{Type: i.Type(), Value: uint64(i.Value)}
}

var integerClassMethods = map[string]RubyMethod{}

var integerMethods = map[string]RubyMethod{
	"div": withArity(1, publicMethod(integerDiv)),
	"/":   withArity(1, publicMethod(integerDiv)),
	"*":   withArity(1, publicMethod(integerMul)),
	"+":   withArity(1, publicMethod(integerAdd)),
	"-":   withArity(1, publicMethod(integerSub)),
	"%":   withArity(1, publicMethod(integerModulo)),
	"<":   withArity(1, publicMethod(integerLt)),
	">":   withArity(1, publicMethod(integerGt)),
	// "==":   withArity(1, publicMethod(integerEq)),
	// "!=":   withArity(1, publicMethod(integerNeq)),
	">=":   withArity(1, publicMethod(integerGte)),
	"<=":   withArity(1, publicMethod(integerLte)),
	"<=>":  withArity(1, publicMethod(integerSpaceship)),
	"to_i": withArity(0, publicMethod(integerToI)),
	"**":   withArity(1, publicMethod(integerPow)),
}

func integerDiv(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	divisor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if divisor.Value == 0 {
		return nil, NewZeroDivisionError()
	}
	return NewInteger(i.Value / divisor.Value), nil
}

func integerMul(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	factor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value * factor.Value), nil
}

func integerAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	add, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value + add.Value), nil
}

func integerSub(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	sub, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value - sub.Value), nil
}

// Objects which can *safely* be converted to an integer
func safeObjectToInteger(arg RubyObject) (int64, bool) {
	var right int64
	switch arg := arg.(type) {
	case *Integer:
		right = int64(arg.Value)
	// case *Boolean:
	// 	if arg.Value {
	// 		right = 1
	// 	} else {
	// 		right = 0
	// 	}
	default:
		return 0, false
	}
	return right, true
}

func integerModulo(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	return NewInteger(i.Value % right), nil
}

func integerLt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value < right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerGt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value > right {
		return TRUE, nil
	}
	return FALSE, nil
}

// func integerEq(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	i := context.Receiver().(*Integer)
// 	right, ok := args[0].(*Integer)
// 	if !ok {
// 		return nil, NewArgumentError(
// 			"comparison of Integer with %s failed",
// 			args[0].Class().(RubyObject).Inspect(),
// 		)
// 	}
// 	if i.Value == right.Value {
// 		return TRUE, nil
// 	}
// 	return FALSE, nil
// }

// func integerNeq(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	i := context.Receiver().(*Integer)
// 	right, ok := args[0].(*Integer)
// 	if !ok {
// 		return nil, NewArgumentError(
// 			"comparison of Integer with %s failed",
// 			args[0].Class().(RubyObject).Inspect(),
// 		)
// 	}
// 	if i.Value != right.Value {
// 		return TRUE, nil
// 	}
// 	return FALSE, nil
// }

func integerSpaceship(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	switch {
	case i.Value > right:
		return &Integer{Value: 1}, nil
	case i.Value < right:
		return &Integer{Value: -1}, nil
	case i.Value == right:
		return &Integer{Value: 0}, nil
	default:
		panic("not reachable")
	}
}

func integerGte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value >= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerLte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	if len(args) != 1 {
		return nil, NewArgumentError("wrong number of arguments (given %d, expected 1)", len(args))
	}
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value <= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerToI(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	return i, nil
}

func integerPow(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	switch arg := args[0].(type) {
	case *Integer:
		result := int64(1)
		for j := int64(0); j < arg.Value; j++ {
			result *= i.Value
		}
		return NewInteger(result), nil
	case *Float:
		result := math.Pow(float64(i.Value), arg.Value)
		return NewFloat(result), nil
	default:
		return nil, NewCoercionTypeError(args[0], i)
	}
}
