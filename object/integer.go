package object

import (
	"fmt"
	"math"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/pkg/errors"
)

var integerClass ruby.ClassObject = newClass("Integer", integerMethods, nil)

func init() {
	CLASSES.Set("Integer", integerClass)
}

// NewInteger returns a new Integer with the given value
//
//go:inline
func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

// Integer represents an integer in Ruby
type Integer struct {
	Value int64
}

// Inspect returns the value as string
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Class returns integerClass
func (i *Integer) Class() ruby.Class { return integerClass }

var (
	_ ruby.Object = &Integer{}
)

func (i *Integer) HashKey() hash.Key {
	return hash.Key(uint64(i.Value))
}

var integerMethods = map[string]ruby.Method{
	"div":  WithArity(1, ruby.NewMethod(integerDiv)),
	"/":    WithArity(1, ruby.NewMethod(integerDiv)),
	"*":    WithArity(1, ruby.NewMethod(integerMul)),
	"+":    WithArity(1, ruby.NewMethod(integerAdd)),
	"-":    WithArity(1, ruby.NewMethod(integerSub)),
	"%":    WithArity(1, ruby.NewMethod(integerModulo)),
	"<":    WithArity(1, ruby.NewMethod(integerLt)),
	">":    WithArity(1, ruby.NewMethod(integerGt)),
	">=":   WithArity(1, ruby.NewMethod(integerGte)),
	"<=":   WithArity(1, ruby.NewMethod(integerLte)),
	"<=>":  WithArity(1, ruby.NewMethod(integerSpaceship)),
	"to_i": WithArity(0, ruby.NewMethod(integerToI)),
	"**":   WithArity(1, ruby.NewMethod(integerPow)),
	"chr":  WithArity(0, ruby.NewMethod(integerChr)),
}

func integerDiv(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	divisor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if divisor.Value == 0 {
		return nil, NewZeroDivisionError()
	}
	return NewInteger(i.Value / divisor.Value), nil
}

func integerMul(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	factor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value * factor.Value), nil
}

func integerAdd(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	add, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value + add.Value), nil
}

func integerSub(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	sub, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value - sub.Value), nil
}

// Objects which can *safely* be converted to an integer
func safeObjectToInteger(arg ruby.Object) (int64, bool) {
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

func integerCmpHelper(args []ruby.Object) (int64, error) {
	right, ok := safeObjectToInteger(args[0])
	if !ok {
		return 0, errors.WithMessage(
			NewArgumentError(
				"comparison of Integer with %s failed",
				args[0].Class().(ruby.Object).Inspect(),
			),
			callersName(),
		)
	}
	return right, nil
}

func integerModulo(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return nil, err
	}
	return NewInteger(i.Value % right), nil
}

func integerLt(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value < right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerGt(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value > right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerSpaceship(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return NIL, err
	}
	switch {
	case i.Value > right:
		return NewInteger(1), nil
	case i.Value < right:
		return NewInteger(-1), nil
	case i.Value == right:
		return NewInteger(0), nil
	default:
		panic("not reachable")
	}
}

func integerGte(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return NIL, err
	}
	if i.Value >= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerLte(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	right, err := integerCmpHelper(args)
	if err != nil {
		return NIL, err
	}
	if i.Value <= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerToI(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	return i, nil
}

func integerPow(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
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

func integerChr(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Integer)
	if i.Value < 0 || i.Value > 255 {
		return nil, NewArgumentError("chr out of range")
	}
	// return NewString(string(i.Value)), nil
	return NewString(string(rune(i.Value))), nil
}
