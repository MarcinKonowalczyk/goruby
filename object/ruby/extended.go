package ruby

type ExtendedObject interface {
	Object
	extendable
	GetMethod(name string) (Method, bool)
	Eigenclass() Eigenclass
}

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	Object
	eigenclass Eigenclass
}

func NewExtendedObject(object Object) ExtendedObject {
	return &extendedObject{
		Object:     object,
		eigenclass: NewEigenclass(object.Class()),
	}
}

func (e *extendedObject) Class() Class    { return e.eigenclass }
func (e *extendedObject) Inspect() string { return e.Object.Inspect() }

func (e *extendedObject) addMethod(name string, method Method) {
	e.eigenclass.AddMethod(name, method)
}

// func (e *extendedObject) GetMethod(
func (e *extendedObject) GetMethod(name string) (Method, bool) {
	if method, ok := e.Object.Class().GetMethod(name); ok {
		return method, true
	}
	return e.eigenclass.GetMethod(name)
}

func (e *extendedObject) Eigenclass() Eigenclass {
	return e.eigenclass
}

var (
	_                Object         = &extendedObject{}
	_                ExtendedObject = &extendedObject{}
	_extended_object ExtendedObject = nil
	_                Object         = _extended_object
)

type extendable interface {
	addMethod(name string, method Method)
}

type extendableRubyObject interface {
	Object
	extendable
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context Object, methodName string, method Method) (Object, bool) {
	extended, is_extendable := context.(extendableRubyObject)
	if !is_extendable {
		extended = NewExtendedObject(context)
	}
	extended.addMethod(methodName, method)
	return extended, !is_extendable
}
