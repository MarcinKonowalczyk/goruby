package object

func newEigenclass(wrappedClass RubyClass, methods map[string]RubyMethod) *eigenclass {
	return &eigenclass{
		methods:      NewMethodSet(methods),
		wrappedClass: wrappedClass,
		Environment:  NewEnvironment(),
	}
}

type eigenclass struct {
	methods      SettableMethodSet
	wrappedClass RubyClass
	Environment
}

func (e *eigenclass) Inspect() string {
	if e.wrappedClass != nil {
		return e.wrappedClass.(RubyClassObject).Inspect()
	}
	return "(singleton class)"
}
func (e *eigenclass) Type() Type { return EIGENCLASS_OBJ }
func (e *eigenclass) Class() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return nil
}
func (e *eigenclass) Methods() MethodSet { return e.methods }
func (e *eigenclass) GetMethod(name string) (RubyMethod, bool) {
	if method, ok := e.methods.Get(name); ok {
		return method, true
	}
	return nil, false
}

func (e *eigenclass) SuperClass() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return bottomClass
}
func (e *eigenclass) New(args ...RubyObject) (RubyObject, error) {
	return e.wrappedClass.New(args...)
}
func (e *eigenclass) Name() string { return e.wrappedClass.Name() }
func (e *eigenclass) addMethod(name string, method RubyMethod) {
	e.methods.Set(name, method)
}
