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

// func methodMissing(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	return nil, NewNoMethodError(context.Receiver(), args[0].(*Symbol).Value)
// }

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context RubyObject, methodName string, method *Function) RubyObject {
	objectToExtend := context
	self, contextIsSelf := context.(*Self)
	if contextIsSelf {
		objectToExtend = self.RubyObject
	}
	extended, contextIsExtendable := objectToExtend.(extendableRubyObject)
	if !contextIsExtendable {
		extended = &extendedObject{
			RubyObject:  objectToExtend,
			class:       newEigenclass(context.Class().(RubyClassObject), map[string]RubyMethod{}),
			Environment: NewEnvironment(),
		}
	}
	extended.addMethod(methodName, method)
	if contextIsSelf {
		self.RubyObject = extended
		return self
	}
	return extended
}
