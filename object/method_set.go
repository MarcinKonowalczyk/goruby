package object

// RubyMethod defines a Ruby method
type RubyMethod interface {
	Call(context CallContext, args ...RubyObject) (RubyObject, error)
}

func withArity(arity int, fn RubyMethod) RubyMethod {
	return &method{
		fn: func(context CallContext, args ...RubyObject) (RubyObject, error) {
			if len(args) != arity {
				return nil, NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(context, args...)
		},
	}
}

func publicMethod(fn func(context CallContext, args ...RubyObject) (RubyObject, error)) RubyMethod {
	return &method{fn: fn}
}

type method struct {
	fn func(context CallContext, args ...RubyObject) (RubyObject, error)
}

func (m *method) Call(context CallContext, args ...RubyObject) (RubyObject, error) {
	return m.fn(context, args...)
}

// MethodSet represents a set of methods
type MethodSet interface {
	// Get returns the method found for name. The boolean will return true if
	// a method was found, false otherwise
	Get(name string) (RubyMethod, bool)
	// GetAll returns a map of name to methods representing the MethodSet.
	GetAll() map[string]RubyMethod
}

// SettableMethodSet represents a MethodSet which can be mutated by setting
// methods on it.
type SettableMethodSet interface {
	MethodSet
	// Set will set method to key name. If there was a method prior defined
	// under name it will be overridden.
	Set(name string, method RubyMethod)
}

// NewMethodSet returns a new method set populated with the given methods
func NewMethodSet(methods map[string]RubyMethod) SettableMethodSet {
	return &methodSet{methods: methods}
}

type methodSet struct {
	methods map[string]RubyMethod
}

func (m *methodSet) GetAll() map[string]RubyMethod {
	methods := make(map[string]RubyMethod)
	for k, v := range m.methods {
		methods[k] = v
	}
	return methods
}

func (m *methodSet) Get(name string) (RubyMethod, bool) {
	method, ok := m.methods[name]
	return method, ok
}

func (m *methodSet) Set(name string, method RubyMethod) {
	m.methods[name] = method
}
