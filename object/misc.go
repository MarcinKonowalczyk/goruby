package object

import (
	"io"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/trace"
)

var (
	CLASSES = env.NewEnvironment[ruby.Object]()
	// symbol to which all function definitions are attached to
	FUNCS_STORE = ruby.NewExtendedObject(FUNCS)
)

// NewMainEnvironment returns a new Environment populated with all Ruby classes
// and the Kernel functions
func NewMainEnvironment() env.Environment[ruby.Object] {
	env := CLASSES.Clone()
	env.SetGlobal("$bottom", BOTTOM)
	env.SetGlobal("$funcs", FUNCS_STORE)
	return env
}

func SetMainEnvironmentIo(env env.Environment[ruby.Object], stdin io.Reader, stdout, stderr io.Writer) {
	if env == nil {
		panic("env cannot be nil")
	}
	stdin_cls := NewIoClass(stdin, nil, nil)
	stdout_cls := NewIoClass(nil, stdout, nil)
	stderr_cls := NewIoClass(nil, nil, stderr)
	env.SetGlobal("$stdin", stdin_cls)
	env.SetGlobal("$stdout", stdout_cls)
	env.SetGlobal("$stderr", stderr_cls)
}

func WithArity(arity int, fn ruby.Method) ruby.Method {
	return ruby.NewMethod(
		func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			defer trace.TraceCtx(ctx, "withArity")()
			if len(args) != arity {
				return nil, NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(ctx, args...)
		},
	)
}

func IsTruthy(obj ruby.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *Integer:
			return obj.Value != 0
		case *Float:
			return obj.Value != 0.0
		case *String:
			return obj.Value != ""
		case *Array:
			return len(obj.Elements) > 0
		case *Hash:
			return len(obj.Map) > 0
		case *Symbol:
			// NOTE: we've checked special symbols above already. other symbols are truthy.
			return true
		default:
			return true
		}
	}
}
