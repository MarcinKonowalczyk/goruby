package env

import (
	"fmt"
	"strings"
)

// NewEnclosedEnvironment returns an Environment wrapped by outer
func NewEnclosedEnvironment[O any](outer Environment[O]) Environment[O] {
	s := make(map[string]O)
	env := &environment[O]{store: s, outer: nil}
	env.outer = outer
	return env
}

// NewEnvironment returns a new Environment ready to use
func NewEnvironment[O any]() Environment[O] {
	s := make(map[string]O)
	return &environment[O]{store: s, outer: nil}
}

// Environment holds Ruby object referenced by strings
type Environment[O any] interface {
	// Get returns the RubyObject found for this key. If it is not found,
	// ok  will be false
	Get(key string) (object O, ok bool)
	// Set sets the RubyObject for the given key. If there is already an
	// object with that key it will be overridden by object
	Set(key string, object O) O
	// Unset removes the entry for the given key. It returns the removed entry
	Unset(key string) O
	// SetGlobal sets val under name at the root of the environment
	SetGlobal(name string, val O) O
	// Outer returns the parent environment
	Outer() Environment[O]
	// Clone returns a copy of the environment. It will shallow copy its values
	//
	// Note that clone will not set its outer env, so calls to Outer will
	// return nil on cloned Environments
	Clone() Environment[O]
	// get something from this environment without checking outer
	// used mainly for testing
	ShallowGet(name string) (O, bool)
}

type environment[O any] struct {
	store map[string]O
	outer Environment[O]
}

// Get returns the RubyObject found for this key. If it is not found,
// ok  will be false
func (e *environment[O]) Get(name string) (O, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *environment[O]) ShallowGet(name string) (O, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set sets the RubyObject for the given key. If there is already an
// object with that key it will be overridden by object
func (e *environment[O]) Set(name string, val O) O {
	e.store[name] = val
	return val
}

// SetGlobal sets val under name at the root of the environment
func (e *environment[O]) SetGlobal(name string, val O) O {
	var env Environment[O] = e
	for env.Outer() != nil {
		env = env.Outer()
	}
	env.Set(name, val)
	return val
}

func (e *environment[O]) Unset(key string) O {
	val := e.store[key]
	delete(e.store, key)
	return val
}

// Outer returns the parent environment
func (e *environment[O]) Outer() Environment[O] {
	return e.outer
}

// Enclose encloses the environment and returns a new one wrapped by outer
func (e *environment[O]) Enclose(outer Environment[O]) Environment[O] {
	env := e.clone()
	env.outer = outer
	return env
}

func (e *environment[O]) Clone() Environment[O] {
	return e.clone()
}

func (e *environment[O]) clone() *environment[O] {
	s := make(map[string]O)
	env := &environment[O]{store: s, outer: nil}
	for k, v := range e.store {
		env.store[k] = v
	}
	return env
}

func (e *environment[O]) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%v", e.store)
	return out.String()
}
