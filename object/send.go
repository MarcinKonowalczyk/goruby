package object

import (
	"hash/fnv"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

// Send sends message method with args to context and returns its result
func Send(ctx call.Context[ruby.Object], method string, args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	trace.MessageCtx(ctx, method)

	receiver := ctx.Receiver()
	receiver_class := receiver.Class()

	// search for the method in the ancestry tree
	for receiver_class != nil {
		fn, ok := receiver_class.GetMethod(method)
		if !ok {
			if rc, ok := receiver_class.(*class); ok && rc == bottomClass {
				// no method and we are at the top of the ancestry tree
				break
			}
			receiver_class = bottomClass
			continue
		}

		return fn.Call(ctx, args...)
	}

	// fmt.Printf("receiver: %v(%T)\n", receiver, receiver)
	// fmt.Printf("method: %v(%T)\n", method, method)
	return nil, NewNoMethodError(receiver, method)
}

func newEigenclass(wrappedClass ruby.Class) *eigenclass {
	return &eigenclass{
		methods:      NewMethodSet(nil),
		wrappedClass: wrappedClass,
	}
}

type eigenclass struct {
	methods      SettableMethodSet
	wrappedClass ruby.Class
}

func (e *eigenclass) Inspect() string {
	if wc, ok := e.wrappedClass.(*class); !ok {
		return e.wrappedClass.(ruby.Object).Inspect()
	} else {
		if wc == nil_class {
			return "(eigenclass of nil)"
		}
		if e.wrappedClass == nil {
			return "(singleton class)"
		}
		return e.wrappedClass.(ruby.Object).Inspect()
	}
}

func (e *eigenclass) Class() ruby.Class {
	return e.wrappedClass
}
func (e *eigenclass) Methods() ruby.MethodSet { return e.methods }
func (e *eigenclass) GetMethod(name string) (ruby.Method, bool) {
	if method, ok := e.methods.Get(name); ok {
		return method, true
	}
	return nil, false
}

func (e *eigenclass) SuperClass() ruby.Class {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return bottomClass
}
func (e *eigenclass) New(args ...ruby.Object) (ruby.Object, error) {
	return e.wrappedClass.New(args...)
}
func (e *eigenclass) Name() string { return e.wrappedClass.Name() }
func (e *eigenclass) addMethod(name string, method ruby.Method) {
	e.methods.Set(name, method)
}
func (e *eigenclass) HashKey() hash.Key {
	if e.wrappedClass != nil {
		h := fnv.New64a()
		h.Write([]byte(e.wrappedClass.Name()))
		return hash.Key(h.Sum64())
	}
	// NOTE: temp fix.
	return hash.Key(1)
}

var (
	_ ruby.Object = &eigenclass{}
	_ ruby.Class  = &eigenclass{}
)

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	ruby.Object
	eigenclass *eigenclass
}

func newExtendedObject(object ruby.Object) *extendedObject {
	object_class := object.Class()
	if object_class == nil {
		panic("object_class is nil")
	}
	if oc, ok := object_class.(*class); ok && oc == nil_class {
		panic("object_class is nil_class")
	}

	return &extendedObject{
		Object:     object,
		eigenclass: newEigenclass(object_class),
	}
}

func (e *extendedObject) Class() ruby.Class { return e.eigenclass }
func (e *extendedObject) Inspect() string   { return e.Object.Inspect() }

// func (e *extendedObject) String() string { return "hello" }
func (e *extendedObject) addMethod(name string, method ruby.Method) {
	e.eigenclass.addMethod(name, method)
}

// func (e *extendedObject) GetMethod(
func (e *extendedObject) GetMethod(name string) (ruby.Method, bool) {
	if method, ok := e.Object.Class().GetMethod(name); ok {
		return method, true
	}
	return e.eigenclass.GetMethod(name)
}

var (
	_ ruby.Object = &extendedObject{}
	_ extendable  = &extendedObject{}
)

type extendable interface {
	addMethod(name string, method ruby.Method)
}

type extendableRubyObject interface {
	ruby.Object
	extendable
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context ruby.Object, methodName string, method *Function) (ruby.Object, bool) {
	extended, is_extendable := context.(extendableRubyObject)
	if !is_extendable {
		extended = newExtendedObject(context)
	}
	extended.addMethod(methodName, method)
	return extended, !is_extendable
}
