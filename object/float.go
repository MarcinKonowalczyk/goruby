package object

import (
	"fmt"
	"unsafe"
)

var floatClass RubyClassObject = newClass(
	"Float", objectClass, floatMethods, floatClassMethods, notInstantiatable,
)

func init() {
	classes.Set("Float", floatClass)
}

// NewFloat returns a new Float with the given value
func NewFloat(value float64) *Float {
	return &Float{Value: value}
}

// Float represents an float in Ruby
type Float struct {
	Value float64
}

// Inspect returns the value as string
func (i *Float) Inspect() string { return fmt.Sprintf("%.16f", i.Value) }

// Type returns FLOAT_OBJ
func (i *Float) Type() Type { return FLOAT_OBJ }

// Class returns floatClass
func (i *Float) Class() RubyClass { return floatClass }

func reinterpretCastFloatToUint64(value float64) uint64 {
	// reinterpret the float as uint64
	value_reinterpret := (*[8]byte)(unsafe.Pointer(&value))[:]
	result := uint64(0)
	for i := 0; i < 8; i++ {
		result |= uint64(value_reinterpret[i]) << (8 * i)
	}
	return result
}
func (i *Float) hashKey() hashKey {
	return hashKey{Type: i.Type(), Value: reinterpretCastFloatToUint64(i.Value)}
}

var floatClassMethods = map[string]RubyMethod{}

var floatMethods = map[string]RubyMethod{
	"div": withArity(1, publicMethod(floatDiv)),
	"/":   withArity(1, publicMethod(floatDiv)),
	"*":   withArity(1, publicMethod(floatMul)),
	"+":   withArity(1, publicMethod(floatAdd)),
	"-":   withArity(1, publicMethod(floatSub)),
	// "%":   withArity(1, publicMethod(floatModulo)),
	"<":    withArity(1, publicMethod(floatLt)),
	">":    withArity(1, publicMethod(floatGt)),
	"==":   withArity(1, publicMethod(floatEq)),
	"!=":   withArity(1, publicMethod(floatNeq)),
	">=":   withArity(1, publicMethod(floatGte)),
	"<=":   withArity(1, publicMethod(floatLte)),
	"<=>":  withArity(1, publicMethod(floatSpaceship)),
	"to_i": withArity(0, publicMethod(floatToI)),
}

func floatDiv(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	var divisor float64
	switch arg := args[0].(type) {
	case *Integer:
		divisor = float64(arg.Value)
	case *Float:
		divisor = arg.Value
	default:
		return nil, NewCoercionTypeError(args[0], i)
	}
	if divisor == 0 {
		return nil, NewZeroDivisionError()
	}
	return NewFloat(i.Value / divisor), nil
}

func floatMul(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	var factor float64
	switch arg := args[0].(type) {
	case *Integer:
		factor = float64(arg.Value)
	case *Float:
		factor = arg.Value
	default:
		return nil, NewCoercionTypeError(args[0], i)
	}

	return NewFloat(i.Value * factor), nil
}

func floatAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	add, ok := args[0].(*Float)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewFloat(i.Value + add.Value), nil
}

func floatSub(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	sub, ok := args[0].(*Float)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewFloat(i.Value - sub.Value), nil
}

// func floatModulo(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	i := context.Receiver().(*Float)
// 	mod, ok := args[0].(*Float)
// 	if !ok {
// 		return nil, NewCoercionTypeError(args[0], i)
// 	}
// 	return NewFloat(i.Value % mod.Value), nil
// }

func floatLt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value < right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatGt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value > right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatEq(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value == right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatNeq(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value != right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatSpaceship(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return NIL, nil
	}
	switch {
	case i.Value > right.Value:
		return &Float{Value: 1}, nil
	case i.Value < right.Value:
		return &Float{Value: -1}, nil
	case i.Value == right.Value:
		return &Float{Value: 0}, nil
	default:
		panic("not reachable")
	}
}

func floatGte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value >= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatLte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	right, ok := args[0].(*Float)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value <= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatToI(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Float)
	return NewInteger(int64(i.Value)), nil
}
