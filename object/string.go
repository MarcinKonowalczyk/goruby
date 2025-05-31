package object

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var stringClass ruby.ClassObject = newClass("String", stringMethods, stringClassMethods)

func init() {
	CLASSES.Set("String", stringClass)
}

//go:inline
func NewString(value string) *String {
	return &String{Value: value}
}

//go:inline
func NewStringf(format string, args ...interface{}) *String {
	return NewString(fmt.Sprintf(format, args...))
}

type String struct {
	Value string
}

func (s *String) Inspect() string   { return s.Value }
func (s *String) Class() ruby.Class { return stringClass }

func (s *String) HashKey() hash.Key {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hash.Key(h.Sum64())
}

var (
	_ ruby.Object = &String{}
)

func stringify(ctx call.Context[ruby.Object], obj ruby.Object) (string, error) {
	if obj == nil {
		return "", NewTypeError(
			"can't convert nil into String",
		)
	}
	ctx2 := call.WithReceiver(ctx, &obj)
	stringObj, err := Send(ctx2, "to_s")
	if err != nil {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s into String",
				obj.Class().Name(),
			),
		)
	}
	str, ok := stringObj.(*String)
	if !ok {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s to String (%s#to_s gives %s)",
				obj.Class().Name(),
				obj.Class().Name(),
				stringObj.Class().Name(),
			),
		)
	}
	return str.Value, nil
}

var stringClassMethods = map[string]ruby.Method{}

var stringMethods = map[string]ruby.Method{
	"to_s":   WithArity(0, ruby.NewMethod(stringToS)),
	"+":      WithArity(1, ruby.NewMethod(stringAdd)),
	"gsub":   WithArity(2, ruby.NewMethod(stringGsub)),
	"length": WithArity(0, ruby.NewMethod(stringLength)),
	"size":   WithArity(0, ruby.NewMethod(stringLength)),
	"lines":  WithArity(0, ruby.NewMethod(stringLines)),
	"to_f":   WithArity(0, ruby.NewMethod(stringToF)),
}

func stringToS(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	str := ctx.Receiver().(*String)
	return NewString(str.Value), nil
}

func stringAdd(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	s := ctx.Receiver().(*String)
	add, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(add, args[0])
	}
	return NewString(s.Value + add.Value), nil
}

func stringGsub(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	s := ctx.Receiver().(*String)
	pattern, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(pattern, args[0])
	}
	replacement, ok := args[1].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(replacement, args[1])
	}

	// Perform the gsub operation
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, NewTypeError(fmt.Sprintf("Invalid regex pattern: %s", err))
	}

	result := re.ReplaceAllString(s.Value, replacement.Value)

	// Return the modified string
	return NewString(result), nil
}

func stringLength(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	s := ctx.Receiver().(*String)
	return NewInteger(int64(len(s.Value))), nil
}

func stringLines(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	s := ctx.Receiver().(*String)
	lines := strings.Split(s.Value, "\n")
	arr := NewArray()
	for _, line := range lines {
		arr.Elements = append(arr.Elements, NewString(line+"\n"))
	}
	return arr, nil
}

var FLOAT_RE = regexp.MustCompile(`[-+]?\d*\.?\d+`)

func stringToF(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	s := ctx.Receiver().(*String)
	if s.Value == "" {
		return NewFloat(0.0), nil
	}
	match := FLOAT_RE.FindString(s.Value)
	if match == "" {
		return NewFloat(0.0), nil
	}
	// Convert the string to a float
	val, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return nil, NewTypeError("Invalid float value: " + s.Value)
	}
	return NewFloat(val), nil
}
