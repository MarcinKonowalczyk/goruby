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

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	RubyObject
	eigenclass *eigenclass
	Environment
}

func (e *extendedObject) Class() RubyClass { return e.eigenclass }
func (e *extendedObject) Inspect() string {
	return e.RubyObject.Inspect()
}

// func (e *extendedObject) String() string { return "hello" }
func (e *extendedObject) addMethod(name string, method RubyMethod) {
	e.eigenclass.addMethod(name, method)
}

type extendable interface {
	addMethod(name string, method RubyMethod)
}

type extendableRubyObject interface {
	RubyObject
	extendable
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context RubyObject, methodName string, method *Function) (RubyObject, bool) {
	objectToExtend := context
	extended, contextIsExtendable := objectToExtend.(extendableRubyObject)
	if !contextIsExtendable {
		extended = &extendedObject{
			RubyObject:  objectToExtend,
			eigenclass:  newEigenclass(context.Class().(RubyClassObject), map[string]RubyMethod{}),
			Environment: NewEnvironment(),
		}
	}
	extended.addMethod(methodName, method)
	return extended, !contextIsExtendable
}
