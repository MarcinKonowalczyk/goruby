package object

import (
	"hash/fnv"

	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
)

var notInstantiatable = func(c ruby.ClassObject, args ...ruby.Object) (ruby.Object, error) {
	return nil, NewNoMethodError(c, "new")
}

// newClass returns a new Ruby Class
func newClass(
	name string,
	instanceMethods,
	classMethods map[string]ruby.Method,
	builder func(ruby.ClassObject, ...ruby.Object) (ruby.Object, error),
) *class {
	if instanceMethods == nil {
		instanceMethods = make(map[string]ruby.Method)
	}
	if classMethods == nil {
		classMethods = make(map[string]ruby.Method)
	}
	cls := &class{
		name:            name,
		instanceMethods: ruby.NewMethodSet(instanceMethods),
		class:           ruby.NewEigenclass(bottomClass),
		builder:         builder,
		Environment:     env.NewEnclosedEnvironment[ruby.Object](nil),
	}
	for name, method := range classMethods {
		cls.class.(ruby.Eigenclass).AddMethod(name, method)
	}
	if cls == nil_class {
		panic("newClass tried to return is nil_class")
	}
	return cls
}

// class represents a Ruby Class object
type class struct {
	name            string
	class           ruby.Class
	instanceMethods ruby.SettableMethodSet
	builder         func(ruby.ClassObject, ...ruby.Object) (ruby.Object, error)
	env.Environment[ruby.Object]
}

func (c *class) Inspect() string         { return c.name }
func (c *class) Class() ruby.Class       { return c.class }
func (c *class) Methods() ruby.MethodSet { return c.instanceMethods }

func (c *class) GetMethod(name string) (ruby.Method, bool) {
	method, ok := c.instanceMethods.Get(name)
	if ok {
		return method, true
	}
	return nil, false
}

func (c *class) HashKey() hash.Key {
	h := fnv.New64a()
	h.Write([]byte(c.name))
	return hash.Key(h.Sum64())
}
func (c *class) New(args ...ruby.Object) (ruby.Object, error) {
	return c.builder(c)
}
func (c *class) Name() string { return c.name }

var (
	_ ruby.Object      = &class{}
	_ ruby.Class       = &class{}
	_ ruby.ClassObject = &class{}
)
