package ruby

import (
	"hash/fnv"

	"github.com/MarcinKonowalczyk/goruby/object/hash"
)

type Eigenclass interface {
	Class
	AddMethod(name string, method Method)
}

func NewEigenclass(wrappedClass Class) Eigenclass {
	return &eigenclass{
		methods:    NewMethodSet(nil),
		superclass: wrappedClass,
	}
}

type eigenclass struct {
	methods    MethodSet
	superclass Class
}

func (e *eigenclass) Inspect() string {
	return e.superclass.(Object).Inspect()
}

func (e *eigenclass) Class() Class {
	return e.superclass
}
func (e *eigenclass) GetMethod(name string) (Method, bool) {
	if method, ok := e.methods.Get(name); ok {
		return method, true
	}
	return nil, false
}

func (e *eigenclass) SuperClass() Class { return e.superclass }
func (e *eigenclass) Name() string      { return e.superclass.Name() }
func (e *eigenclass) AddMethod(name string, method Method) {
	e.methods.Set(name, method)
}
func (e *eigenclass) HashKey() hash.Key {
	if e.superclass != nil {
		h := fnv.New64a()
		h.Write([]byte(e.superclass.Name()))
		return hash.Key(h.Sum64())
	}
	// NOTE: temp fix.
	return hash.Key(1)
}

var (
	_ Eigenclass = &eigenclass{}
	_ Object     = &eigenclass{}
	_ Class      = &eigenclass{}
)
