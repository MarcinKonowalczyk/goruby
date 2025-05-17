package object

import (
	"hash/fnv"
)

var notInstantiatable = func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
	return nil, NewNoMethodError(c, "new")
}

// newClass returns a new Ruby Class
func newClass(
	name string,
	instanceMethods,
	classMethods map[string]RubyMethod,
	builder func(RubyClassObject, ...RubyObject) (RubyObject, error),
) *class {
	if instanceMethods == nil {
		instanceMethods = make(map[string]RubyMethod)
	}
	if classMethods == nil {
		classMethods = make(map[string]RubyMethod)
	}
	cls := &class{
		name:            name,
		instanceMethods: NewMethodSet(instanceMethods),
		class:           newEigenclass(bottomClass),
		builder:         builder,
		Environment:     NewEnclosedEnvironment(nil),
	}
	cls.class.(*eigenclass).methods = NewMethodSet(classMethods)
	return cls
}

// class represents a Ruby Class object
type class struct {
	name            string
	class           RubyClass
	instanceMethods SettableMethodSet
	builder         func(RubyClassObject, ...RubyObject) (RubyObject, error)
	Environment
}

func (c *class) Inspect() string    { return c.name }
func (c *class) Class() RubyClass   { return c.class }
func (c *class) Methods() MethodSet { return c.instanceMethods }

func (c *class) GetMethod(name string) (RubyMethod, bool) {
	method, ok := c.instanceMethods.Get(name)
	if ok {
		return method, true
	}
	return nil, false
}

func (c *class) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(c.name))
	return HashKey(h.Sum64())
}
func (c *class) addMethod(name string, method RubyMethod) {
	c.instanceMethods.Set(name, method)
}
func (c *class) New(args ...RubyObject) (RubyObject, error) {
	return c.builder(c)
}
func (c *class) Name() string { return c.name }

var (
	_ RubyObject = &class{}
	_ RubyClass  = &class{}
)
