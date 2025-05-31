package object

import (
	"fmt"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var (
	nil_class   *class              = nil
	bottomClass *class              = nil_class
	BOTTOM      ruby.ExtendedObject = nil
)

// TODO: make sure we don't collide with other hash keys
const HASH_KEY_BOTTOM hash.Key = 0

func init() {
	// NOTE: create the bottom class in init to avoid circular import
	bottomClass = newClass("Bottom", bottomMethodSet, nil)
	bottomClass.class = bottomClass
	CLASSES.Set("Bottom", bottomClass)

	BOTTOM = ruby.NewExtendedObject(&Bottom{})
}

type Bottom struct{}

func (o *Bottom) Inspect() string   { return "" }
func (o *Bottom) Class() ruby.Class { return bottomClass }
func (o *Bottom) HashKey() hash.Key { return HASH_KEY_BOTTOM }

var bottomMethodSet = map[string]ruby.Method{
	"to_s":    WithArity(0, ruby.NewMethod(bottomToS)),
	"is_a?":   WithArity(1, ruby.NewMethod(bottomIsA)),
	"nil?":    WithArity(0, ruby.NewMethod(bottomIsNil)),
	"methods": ruby.NewMethod(bottomMethods),
	"class":   WithArity(0, ruby.NewMethod(bottomClassMethod)),
	"puts":    ruby.NewMethod(bottomPuts),
	"print":   ruby.NewMethod(bottomPrint),
	"raise":   ruby.NewMethod(bottomRaise),
	"==":      WithArity(1, ruby.NewMethod(bottomEqual)),
	"!=":      WithArity(1, ruby.NewMethod(bottomNotEqual)),
}

func bottomToS(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	receiver := ctx.Receiver()
	val := fmt.Sprintf("#<%s:%p>", receiver.Class().Name(), receiver)
	return NewString(val), nil
}

func bottomIsA(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	receiver_class := ctx.Receiver().Class()
	switch arg := args[0].(type) {
	case ruby.ClassObject:
		if arg.Name() == receiver_class.Name() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return nil, NewTypeError("argument must be a Class")
	}
}

func print(lines []string, delimiter string) {
	var out strings.Builder
	for i, line := range lines {
		out.WriteString(line)
		if i != len(lines)-1 {
			out.WriteString(delimiter)
		}
	}
	out.WriteString(delimiter)
	fmt.Print(out.String())
}

func bottomPuts(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	var lines []string
	for _, arg := range args {
		if arr, ok := arg.(*Array); ok {
			// arg is an array. splat it out
			// todo: make it a deep splat? check with original ruby implementation
			for _, elem := range arr.Elements {
				lines = append(lines, elem.Inspect())
			}
		} else {
			switch arg := arg.(type) {
			case *Symbol:
				if arg == NIL.(*Symbol) {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	print(lines, "\n")
	return NIL, nil
}

func bottomPrint(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	var lines []string
	for _, arg := range args {
		if arr, ok := arg.(*Array); ok {
			// arg is an array. splat it out
			// todo: make it a deep splat? check with original ruby implementation
			// for _, elem := range arr.Elements {
			// 	lines = append(lines, elem.Inspect())
			// }
			lines = append(lines, arr.Inspect())
		} else {
			switch arg := arg.(type) {
			case *Symbol:
				if arg == NIL.(*Symbol) {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	print(lines, "")
	return NIL, nil
}

func bottomMethods(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	showSuperMethods := true
	if len(args) == 1 {
		if val, ok := SymbolToBool(args[0]); ok {
			showSuperMethods = val
		}
	}

	receiver := ctx.Receiver()
	class := ctx.Receiver().Class()

	extended, ok := receiver.(ruby.ExtendedObject)

	if !showSuperMethods && !ok {
		return &Array{}, nil
	}

	if !showSuperMethods && ok {
		class = extended.Eigenclass()
	}

	return getMethods(class, showSuperMethods), nil
}

func bottomIsNil(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	receiver := ctx.Receiver()
	if receiver == NIL {
		return TRUE, nil
	}
	return FALSE, nil
}

func bottomClassMethod(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	receiver := ctx.Receiver()
	if _, ok := receiver.(ruby.ClassObject); ok {
		return nil, nil
	}
	return receiver.Class().(ruby.ClassObject), nil
}

func bottomRaise(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	switch len(args) {
	case 1:
		switch arg := args[0].(type) {
		case *String:
			return nil, NewRuntimeError("%s", arg.Value)
		default:
			return nil, NewRuntimeError("%s", arg.Inspect())
			// default:
			// 	exc, err := Send(NewCallContext(context.Env(), arg), "exception")
			// 	if err != nil {
			// 		return nil, NewTypeError("exception class/object expected")
			// 	}
			// 	if excAsErr, ok := exc.(error); ok {
			// 		return nil, excAsErr
			// 	}
			// 	return nil, nil
		}
	default:
		return nil, NewRuntimeError("")
	}
}

func bottomEqual(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	if RubyObjectsEqual(ctx.Receiver(), args[0]) {
		return TRUE, nil
	}
	return FALSE, nil
}

func bottomNotEqual(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	if RubyObjectsEqual(ctx.Receiver(), args[0]) {
		return FALSE, nil
	}
	return TRUE, nil
}
