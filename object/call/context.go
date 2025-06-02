package call

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/object/env"
)

type _R interface{} // for type assertions

////////////////////////////////////////////////////////////////////////////////

type Io interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type EvalContext[R any] interface {
	Receiver() R
	// Env returns the current environment at call time
	Env() env.Environment[R]
	// Eval represents an evaluation method suitable to eval arbitrary Ruby
	// AST Nodes and transform them into a resulting Ruby object or an error.
	Eval(ast.Node, env.Environment[R]) (R, error)
}

// carries all the information about a call + all the context.Context stuff
// like loggers, tracers, deadlines, etc.
type Context[R any] interface {
	context.Context
	EvalContext[R]
	Io

	WithReceiver(receiver R) Context[R]
	WithEnv(env env.Environment[R]) Context[R]
	WithEval(eval EvalFunc[R]) Context[R]
}

type EvalFunc[R any] func(ast.Node, env.Environment[R]) (R, error)

// NewContext returns a new CallContext with a stubbed eval function
func NewContext[R any](ctx context.Context, io Io) Context[R] {
	return emptyContext[R]{ctx, io}
}

////////////////////////////////////////////////////////////////////////////////

// utility struct to wrap a context and replace the receiver, env or eval
type wrappedCtx[R any] struct {
	Context[R]
	io       Io
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

func (c *wrappedCtx[R]) Io() Io {
	if c.io != nil {
		return c.io
	}
	return c.Context.(Io)
}

//go:inline
func panicIfNil[T any](name string, value *T) {
	if value == nil {
		panic(fmt.Sprintf("%s cannot be nil", name))
	}
}

func (c wrappedCtx[R]) WithReceiver(receiver R) Context[R] {
	panicIfNil("receiver", &receiver)
	c.receiver = &receiver
	return &c
}

func (c wrappedCtx[R]) WithEnv(env env.Environment[R]) Context[R] {
	panicIfNil("env", &env)
	c.env = &env
	return &c
}

func (c wrappedCtx[R]) WithEval(eval EvalFunc[R]) Context[R] {
	panicIfNil("eval", &eval)
	c.eval = &eval
	return &c
}

var (
	_ Context[_R]     = (*wrappedCtx[_R])(nil)
	_ context.Context = (*wrappedCtx[_R])(nil)
)

// func WithReceiver[R any](parent Context[R], receiver *R) Context[R] {
// 	panicIfNil("parent", &parent)
// 	panicIfNil("parent", &parent)
// 	return &wrappedCtx[R]{Context: parent, receiver: receiver, env: nil, eval: nil}
// }

// func WithEnv[R any](parent Context[R], env *env.Environment[R]) Context[R] {
// 	panicIfNil("parent", &parent)
// 	panicIfNil("env", &env)
// 	return &wrappedCtx[R]{Context: parent, receiver: nil, env: env, eval: nil}
// }

// func WithEval[R any](parent Context[R], eval EvalFunc[R]) Context[R] {
// 	panicIfNil("parent", &parent)
// 	return &wrappedCtx[R]{Context: parent, receiver: nil, env: nil, eval: &eval}
// }

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

type emptyContext[R any] struct {
	context.Context // embedded context.Context
	Io
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

func (c emptyContext[R]) Stdin() io.Reader {
	return os.Stdin
}

func (c emptyContext[R]) Stdout() io.Writer {
	return os.Stdout
}

func (c emptyContext[R]) Stderr() io.Writer {
	return os.Stderr
}

func (c emptyContext[R]) WithReceiver(receiver R) Context[R] {
	panicIfNil("receiver", &receiver)
	return &wrappedCtx[R]{Context: c, receiver: &receiver, env: nil, eval: nil}
}

func (c emptyContext[R]) WithEnv(env env.Environment[R]) Context[R] {
	panicIfNil("env", &env)
	return &wrappedCtx[R]{Context: c, receiver: nil, env: &env, eval: nil}
}

func (c emptyContext[R]) WithEval(eval EvalFunc[R]) Context[R] {
	panicIfNil("eval", &eval)
	return &wrappedCtx[R]{Context: c, receiver: nil, env: nil, eval: &eval}
}

var (
	_ Context[_R]     = (*emptyContext[_R])(nil)
	_ context.Context = (*emptyContext[_R])(nil)
	_ Io              = (*emptyContext[_R])(nil)
)
