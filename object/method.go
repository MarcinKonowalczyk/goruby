package object

import (
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

// type CC call.Context[ruby.Object]

// type ruby.Method interface {
// 	Call(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error)
// }

func withArity(arity int, fn ruby.Method) ruby.Method {
	return &method{
		fn: func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			defer trace.TraceCtx(ctx, "withArity")()
			if len(args) != arity {
				return nil, NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(ctx, args...)
		},
	}
}

func newMethod(fn func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error)) ruby.Method {
	return &method{fn: fn}
}

type method struct {
	fn func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error)
}

func (m *method) Call(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, "method.Call")()
	return m.fn(ctx, args...)
}

// SettableMethodSet represents a MethodSet which can be mutated by setting
// methods on it.
type SettableMethodSet interface {
	ruby.MethodSet
	// Set will set method to key name. If there was a method prior defined
	// under name it will be overridden.
	Set(name string, method ruby.Method)
}

// NewMethodSet returns a new method set populated with the given methods
func NewMethodSet(methods map[string]ruby.Method) SettableMethodSet {
	if methods == nil {
		methods = make(map[string]ruby.Method)
	}
	return &methodSet{methods: methods}
}

type methodSet struct {
	methods map[string]ruby.Method
}

func (m *methodSet) Names() []string {
	methods := make([]string, 0, len(m.methods))
	for k := range m.methods {
		methods = append(methods, k)
	}
	return methods
}

func (m *methodSet) Get(name string) (ruby.Method, bool) {
	method, ok := m.methods[name]
	return method, ok
}

func (m *methodSet) Set(name string, method ruby.Method) {
	m.methods[name] = method
}
