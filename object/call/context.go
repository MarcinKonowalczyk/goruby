package call

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/object/env"
)

type _R interface{} // for type assertions

////////////////////////////////////////////////////////////////////////////////

// carries all the information about a call + all the context.Context stuff
// like loggers, tracers, deadlines, etc.
// R and E are the types of the receiver and the environment. they needed
// to be generics to break the import cycle with the object package, but
// in practise they are always RubyObject and Environment respectively.
type Context[R any] interface {
	// embedded context.Context
	context.Context
	// Receiver returns the Ruby object the message is sent to
	Receiver() R
	// Env returns the current environment at call time
	Env() env.Environment[R]
	// Eval represents an evaluation method suitable to eval arbitrary Ruby
	// AST Nodes and transform them into a resulting Ruby object or an error.
	Eval(ast.Node, env.Environment[R]) (R, error)
}

type EvalFunc[R any] func(ast.Node, env.Environment[R]) (R, error)

// NewContext returns a new CallContext with a stubbed eval function
func NewContext[R any](receiver R, env env.Environment[R]) Context[R] {
	return WrappedContext(
		emptyContext[R]{},
		&receiver,
		&env,
		nil,
	)
}

////////////////////////////////////////////////////////////////////////////////

// utility struct to wrap a context and replace the receiver, env or eval
type wrappedCtx[R any] struct {
	Context[R]
	receiver *R
	env      *env.Environment[R]
	eval     *EvalFunc[R]
}

func (c *wrappedCtx[R]) Receiver() R {
	if c.receiver != nil {
		return *c.receiver
	}
	return c.Context.Receiver()
}

func (c *wrappedCtx[R]) Env() env.Environment[R] {
	if c.env != nil {
		return *c.env
	}
	return c.Context.Env()
}

func (c *wrappedCtx[R]) Eval(node ast.Node, env env.Environment[R]) (R, error) {
	if c.eval != nil {
		return (*c.eval)(node, env)
	}
	return c.Context.Eval(node, env)
}

var (
	_ Context[_R]     = (*wrappedCtx[_R])(nil)
	_ context.Context = (*wrappedCtx[_R])(nil)
)

func WithReceiver[R any](parent Context[R], receiver *R) Context[R] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if receiver == nil {
		panic("receiver cannot be nil")
	}
	return &wrappedCtx[R]{
		Context:  parent,
		receiver: receiver,
		env:      nil,
		eval:     nil,
	}
}

func WithEnv[R any](parent Context[R], env *env.Environment[R]) Context[R] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if env == nil {
		panic("env cannot be nil")
	}
	return &wrappedCtx[R]{
		Context:  parent,
		receiver: nil,
		env:      env,
		eval:     nil,
	}
}

func WithEval[R any](parent Context[R], eval EvalFunc[R]) Context[R] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if eval == nil {
		panic("eval function cannot be nil")
	}
	return &wrappedCtx[R]{
		Context:  parent,
		receiver: nil,
		env:      nil,
		eval:     &eval,
	}
}

func WrappedContext[R any](
	parent Context[R],
	receiver *R,
	env *env.Environment[R],
	eval EvalFunc[R],
) Context[R] {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	if receiver == nil && env == nil && eval == nil {
		panic("at least one of receiver, env or eval must be non-nil")
	}
	return &wrappedCtx[R]{
		Context:  parent,
		receiver: receiver,
		env:      env,
		eval:     &eval,
	}
}

////////////////////////////////////////////////////////////////////////////////

type emptyContext[R any] struct{}

func (emptyContext[R]) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyContext[R]) Done() <-chan struct{} {
	return nil
}

func (emptyContext[R]) Err() error {
	return nil
}

func (emptyContext[R]) Value(key any) any {
	return nil
}

func (c emptyContext[R]) Receiver() R {
	var zero R
	return zero
}

func (c emptyContext[R]) Env() env.Environment[R] {
	var zero env.Environment[R]
	return zero
}

func (c emptyContext[R]) Eval(node ast.Node, env env.Environment[R]) (R, error) {
	var zero R
	return zero, fmt.Errorf("no eval present")
}

var (
	_ Context[_R]     = (*emptyContext[_R])(nil)
	_ context.Context = (*emptyContext[_R])(nil)
)
