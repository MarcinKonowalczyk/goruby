package ruby

import (
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

func NewMethod(fn func(ctx call.Context[Object], args ...Object) (Object, error)) Method {
	return &method{fn: fn}
}

type method struct {
	fn func(ctx call.Context[Object], args ...Object) (Object, error)
}

func (m *method) Call(ctx call.Context[Object], args ...Object) (Object, error) {
	defer trace.TraceCtx(ctx, "method.Call")()
	return m.fn(ctx, args...)
}

// SettableMethodSet represents a MethodSet which can be mutated by setting
// methods on it.
type SettableMethodSet interface {
	MethodSet
	// Set will set method to key name. If there was a method prior defined
	// under name it will be overridden.
	Set(name string, method Method)
}

// NewMethodSet returns a new method set populated with the given methods
func NewMethodSet(methods map[string]Method) SettableMethodSet {
	if methods == nil {
		methods = make(map[string]Method)
	}
	return &methodSet{methods: methods}
}

type methodSet struct {
	methods map[string]Method
}

func (m *methodSet) Names() []string {
	methods := make([]string, 0, len(m.methods))
	for k := range m.methods {
		methods = append(methods, k)
	}
	return methods
}

func (m *methodSet) Get(name string) (Method, bool) {
	method, ok := m.methods[name]
	return method, ok
}

func (m *methodSet) Set(name string, method Method) {
	m.methods[name] = method
}
