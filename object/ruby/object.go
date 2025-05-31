package ruby

import (
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
)

type inspectable interface {
	Inspect() string
}

type hashable interface {
	HashKey() hash.Key
}

// Object represents an object in Ruby
type Object interface {
	inspectable
	hashable
	Class() Class
}

type Method interface {
	Call(ctx call.Context[Object], args ...Object) (Object, error)
}

// MethodSet represents a set of methods
type MethodSet interface {
	// Get returns the method found for name. The boolean will return true if
	// a method was found, false otherwise
	Get(name string) (Method, bool)
	// Names returns the names of all methods in the set
	Names() []string
}

// Class represents a class in Ruby
type Class interface {
	inspectable
	hashable
	GetMethod(name string) (Method, bool)
	Methods() MethodSet
	New(args ...Object) (Object, error)
	Name() string
}

// ClassObject represents a class object in Ruby
type ClassObject interface {
	Object
	Class
}
