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

// Class represents a class in Ruby
type Class interface {
	inspectable
	hashable
	GetMethod(name string) (Method, bool)
	Name() string
}

// ClassObject represents a class object in Ruby
type ClassObject interface {
	Object
	Class
}
