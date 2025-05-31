package call

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

// types for type assertions
type _R interface{}
type _E interface{}

////////////////////////////////////////////////////////////////////////////////

// carries all the information about a call + all the context.Context stuff
// like loggers, tracers, deadlines, etc.
// R and E are the types of the receiver and the environment. they needed
// to be generics to break the import cycle with the object package, but
// in practise they are always RubyObject and Environment respectively.
type Context[R, E any] interface {
	// embedded context.Context
	context.Context
	// Receiver returns the Ruby object the message is sent to
	Receiver() R
	// Env returns the current environment at call time
	Env() E
	// Eval represents an evaluation method suitable to eval arbitrary Ruby
	// AST Nodes and transform them into a resulting Ruby object or an error.
	Eval(ast.Node, E) (R, error)
}

// NewContext returns a new CallContext with a stubbed eval function
func NewContext[R, E any](receiver R, env E) Context[R, E] {
	return WrappedContext(
		emptyContext[R, E]{},
		&receiver,
		&env,
		nil,
	)
}

////////////////////////////////////////////////////////////////////////////////

// utility struct to wrap a context and replace the receiver, env or eval
type wrappedCtx[R, E any] struct {
	Context[R, E]
	receiver *R
	env      *E
	eval     *func(ast.Node, E) (R, error)
}

func (c *wrappedCtx[R, E]) Receiver() R {
	if c.receiver != nil {
		return *c.receiver
	}
	return c.Context.Receiver()
}

func (c *wrappedCtx[R, E]) Env() E {
	if c.env != nil {
		return *c.env
	}
	return c.Context.Env()
}

func (c *wrappedCtx[R, E]) Eval(node ast.Node, env E) (R, error) {
	if c.eval != nil {
		return (*c.eval)(node, env)
	}
	return c.Context.Eval(node, env)
}

var (
	_ Context[_R, _E] = (*wrappedCtx[_R, _E])(nil)
	_ context.Context = (*wrappedCtx[_R, _E])(nil)
)

func WithReceiver[R, E any](parent Context[R, E], receiver *R) Context[R, E] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if receiver == nil {
		panic("receiver cannot be nil")
	}
	return &wrappedCtx[R, E]{
		Context:  parent,
		receiver: receiver,
		env:      nil,
		eval:     nil,
	}
}

func WithEnv[R, E any](parent Context[R, E], env *E) Context[R, E] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if env == nil {
		panic("env cannot be nil")
	}
	return &wrappedCtx[R, E]{
		Context:  parent,
		receiver: nil,
		env:      env,
		eval:     nil,
	}
}

func WithEval[R, E any](parent Context[R, E], eval func(ast.Node, E) (R, error)) Context[R, E] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if eval == nil {
		panic("eval function cannot be nil")
	}
	return &wrappedCtx[R, E]{
		Context:  parent,
		receiver: nil,
		env:      nil,
		eval:     &eval,
	}
}

func WrappedContext[R, E any](parent Context[R, E], receiver *R, env *E, eval func(ast.Node, E) (R, error)) Context[R, E] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if receiver == nil && env == nil && eval == nil {
		panic("at least one of receiver, env or eval must be non-nil")
	}
	return &wrappedCtx[R, E]{
		Context:  parent,
		receiver: receiver,
		env:      env,
		eval:     &eval,
	}
}

////////////////////////////////////////////////////////////////////////////////

type emptyContext[R, E any] struct{}

func (emptyContext[R, E]) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyContext[R, E]) Done() <-chan struct{} {
	return nil
}

func (emptyContext[R, E]) Err() error {
	return nil
}

func (emptyContext[R, E]) Value(key any) any {
	return nil
}

func (c emptyContext[R, E]) Receiver() R {
	var zero R
	return zero
}

func (c emptyContext[R, E]) Env() E {
	var zero E
	return zero
}

func (c emptyContext[R, E]) Eval(node ast.Node, env E) (R, error) {
	var zero R
	return zero, fmt.Errorf("no eval present")
}

var (
	_ Context[_R, _E] = (*emptyContext[_R, _E])(nil)
	_ context.Context = (*emptyContext[_R, _E])(nil)
)
