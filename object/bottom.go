package object

import (
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/trace"
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
	"to_s":  WithArity(0, ruby.NewMethod(bottomToS)),
	"is_a?": WithArity(1, ruby.NewMethod(bottomIsA)),
	"nil?":  WithArity(0, ruby.NewMethod(bottomIsNil)),
	"class": WithArity(0, ruby.NewMethod(bottomClassMethod)),
	"puts":  ruby.NewMethod(bottomPuts),
	"print": ruby.NewMethod(bottomPrint),
	"raise": ruby.NewMethod(bottomRaise),
	"==":    WithArity(1, ruby.NewMethod(bottomEqual)),
	"!=":    WithArity(1, ruby.NewMethod(bottomNotEqual)),
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

func bottomPuts(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	stdout, ok := ctx.Env().Get("$stdout")
	if !ok {
		return nil, NewRuntimeError("no $stdout defined in the environment")
	}
	return Send(ctx.WithReceiver(stdout), "puts", args...)
}

func bottomPrint(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	stdout, ok := ctx.Env().Get("$stdout")
	if !ok {
		return nil, NewRuntimeError("no $stdout defined in the environment")
	}
	return Send(ctx.WithReceiver(stdout), "print", args...)
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
