package object

// Send sends message method with args to context and returns its result
func Send(context CallContext, method string, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	class := receiver.Class()

	// search for the method in the ancestry tree
	for class != nil {
		fn, ok := class.GetMethod(method)
		if !ok {
			if class == bottomClass {
				// no method and we are at the top of the ancestry tree
				break
			}
			class = bottomClass
			continue
		}

		return fn.Call(context, args...)
	}

	return nil, NewNoMethodError(receiver, method)
}

func newEigenclass(wrappedClass RubyClass) *eigenclass {
	return &eigenclass{
		methods:      NewMethodSet(nil),
		wrappedClass: wrappedClass,
	}
}

type eigenclass struct {
	methods      SettableMethodSet
	wrappedClass RubyClass
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

var (
	_ RubyObject = &eigenclass{}
	_ RubyClass  = &eigenclass{}
)

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	RubyObject
	eigenclass *eigenclass
}

func newExtendedObject(object RubyObject) *extendedObject {
	return &extendedObject{
		RubyObject: object,
		eigenclass: newEigenclass(object.Class()),
	}
}

func (e *extendedObject) Class() RubyClass { return e.eigenclass }
func (e *extendedObject) Inspect() string {
	return e.RubyObject.Inspect()
}

// func (e *extendedObject) String() string { return "hello" }
func (e *extendedObject) addMethod(name string, method RubyMethod) {
	e.eigenclass.addMethod(name, method)
}

var (
	_ RubyObject = &extendedObject{}
	_ extendable = &extendedObject{}
)

type extendable interface {
	addMethod(name string, method RubyMethod)
}

type extendableRubyObject interface {
	RubyObject
	extendable
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context RubyObject, methodName string, method *Function) (RubyObject, bool) {
	extended, is_extendable := context.(extendableRubyObject)
	if !is_extendable {
		extended = newExtendedObject(context)
	}
	extended.addMethod(methodName, method)
	return extended, !is_extendable
}
