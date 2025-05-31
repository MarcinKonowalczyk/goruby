package call

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

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

func noEvalPresent[R, E any](node ast.Node, env E) (R, error) {
	return *new(R), fmt.Errorf("no eval present")
}

// NewContext returns a new CallContext with a stubbed eval function
func NewContext[R, E any](receiver R, env E) Context[R, E] {
	return &callContext[R, E]{
		emptyCtx: emptyCtx{},
		receiver: receiver,
		env:      env,
		eval:     noEvalPresent[R, E],
	}
}

// func NewWithReceiver(parent context.Context, receiver RubyObject) CallContext {
// 	if parent == nil {
// 		panic("parent context cannot be nil")
// 	}
// 	if ctx, ok := parent.(CallContext); ok {
// 		return WithReceiver(ctx, receiver)
// 	}
// 	return &callContext{
// 		emptyCtx: emptyCtx{},
// 		env:      NewMainEnvironment(),
// 		eval:     noEvalPresent,
// 		receiver: receiver,
// 	}
// }

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

// types for type assertions
type _R interface{}
type _E interface{}

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

func WrappedContext[R, E any](parent Context[R, E], receiver *R, env *E, eval *func(ast.Node, E) (R, error)) Context[R, E] {
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
		eval:     eval,
	}
}

////////////////////////////////////////////////////////////////////////////////

type emptyCtx struct{}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyCtx) Done() <-chan struct{} {
	return nil
}

func (emptyCtx) Err() error {
	return nil
}

func (emptyCtx) Value(key any) any {
	return nil
}

type callContext[R, E any] struct {
	emptyCtx
	receiver R
	env      E
	eval     func(node ast.Node, env E) (R, error)
}

func (c *callContext[R, E]) Env() E { return c.env }
func (c *callContext[R, E]) Eval(node ast.Node, env E) (R, error) {
	return c.eval(node, env)
}
func (c *callContext[R, E]) Receiver() R { return c.receiver }

var (
	_ Context[_R, _E] = (*callContext[_R, _E])(nil)
	_ context.Context = (*callContext[_R, _E])(nil)
)
