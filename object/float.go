package object

import (
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/pkg/errors"
)

var floatClass ruby.ClassObject = newClass(
	"Float", floatMethods, nil, notInstantiatable,
)

func init() {
	CLASSES.Set("Float", floatClass)
}

// NewFloat returns a new Float with the given value
func NewFloat(value float64) *Float {
	return &Float{Value: value}
}

// Float represents an float in Ruby
type Float struct {
	Value float64
}

func (i *Float) Inspect() string   { return fmt.Sprintf("%.16f", i.Value) }
func (i *Float) Class() ruby.Class { return floatClass }

func reinterpretCastFloatToUint64(value float64) uint64 {
	// reinterpret the float as uint64
	value_reinterpret := (*[8]byte)(unsafe.Pointer(&value))[:]
	result := uint64(0)
	for i := 0; i < 8; i++ {
		result |= uint64(value_reinterpret[i]) << (8 * i)
	}
	return result
}

func (i *Float) HashKey() hash.Key {
	return hash.Key(reinterpretCastFloatToUint64(i.Value))
}

var (
	_ ruby.Object = &Float{}
)

var floatMethods = map[string]ruby.Method{
	"div":  WithArity(1, ruby.NewMethod(floatDiv)),
	"/":    WithArity(1, ruby.NewMethod(floatDiv)),
	"*":    WithArity(1, ruby.NewMethod(floatMul)),
	"+":    WithArity(1, ruby.NewMethod(floatAdd)),
	"-":    WithArity(1, ruby.NewMethod(floatSub)),
	"<":    WithArity(1, ruby.NewMethod(floatLt)),
	">":    WithArity(1, ruby.NewMethod(floatGt)),
	">=":   WithArity(1, ruby.NewMethod(floatGte)),
	"<=":   WithArity(1, ruby.NewMethod(floatLte)),
	"<=>":  WithArity(1, ruby.NewMethod(floatSpaceship)),
	"to_i": WithArity(0, ruby.NewMethod(floatToI)),
	"**":   WithArity(1, ruby.NewMethod(floatPow)),
}

func floatDiv(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
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

func floatMul(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
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

func floatAdd(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	add, ok := args[0].(*Float)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewFloat(i.Value + add.Value), nil
}

func floatSub(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	sub, ok := args[0].(*Float)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewFloat(i.Value - sub.Value), nil
}

// Objects which can *safely* be converted to a float
func safeObjectToFloat(arg ruby.Object) (float64, bool) {
	var right float64
	switch arg := arg.(type) {
	case *Float:
		right = arg.Value
	case *Integer:
		right = float64(arg.Value)
	default:
		return 0, false
	}
	return right, true
}

func callersName() string {
	parent, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(parent)
	return fn.Name()
}

func floatCmpHelper(args []ruby.Object) (float64, error) {
	if len(args) != 1 {
		return 0, errors.WithMessage(
			NewWrongNumberOfArgumentsError(1, len(args)),
			callersName(),
		)
	}
	right, ok := safeObjectToFloat(args[0])
	if !ok {
		return 0, errors.WithMessage(
			NewArgumentError(
				"comparison of Float with %s failed",
				args[0].Class().(ruby.Object).Inspect(),
			),
			callersName(),
		)
	}
	return right, nil
}
func floatLt(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	right, err := floatCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value < right {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatGt(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	right, err := floatCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value > right {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatSpaceship(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	right, err := floatCmpHelper(args)
	if err != nil {
		return NIL, err
	}
	switch {
	case i.Value > right:
		return NewFloat(1), nil
	case i.Value < right:
		return NewFloat(-1), nil
	case i.Value == right:
		return NewFloat(0), nil
	default:
		panic("not reachable")
	}
}

func floatGte(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	right, err := floatCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value >= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatLte(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	right, err := floatCmpHelper(args)
	if err != nil {
		return nil, err
	}
	if i.Value <= right {
		return TRUE, nil
	}
	return FALSE, nil
}

func floatToI(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	return NewInteger(int64(i.Value)), nil
}

func floatPow(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	i := ctx.Receiver().(*Float)
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	right, ok := safeObjectToFloat(args[0])
	if !ok {
		return nil, NewArgumentError(
			"comparison of Float with %s failed",
			args[0].Class().(ruby.Object).Inspect(),
		)
	}
	if right < 0 {
		return nil, NewArgumentError("negative exponent")
	}
	result := math.Pow(i.Value, right)
	return NewFloat(result), nil
}
