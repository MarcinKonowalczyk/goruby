package object

import (
	"context"
	"fmt"
	"time"

	"github.com/MarcinKonowalczyk/goruby/ast"
)

// CallContext represents the context information when sending a message to an
// object.
type CallContext interface {
	// embedded context.Context
	context.Context
	// Env returns the current environment at call time
	Env() Environment
	// Eval represents an evalualtion method suitable to eval arbitrary Ruby
	// AST Nodes and transform them into a resulting Ruby object or an error.
	Eval(ast.Node, Environment) (RubyObject, error)
	// EvalCtx(ast.Node, Environment) (RubyObject, error)
	// Receiver returns the Ruby object the message is sent to
	Receiver() RubyObject
}

var noEvalPresent = func(node ast.Node, env Environment) (RubyObject, error) {
	return nil, fmt.Errorf("no eval present")
}

// NewCallContext returns a new CallContext with a stubbed eval function
func NewCallContext(env Environment, receiver RubyObject) CallContext {
	return &callContext{
		emptyCtx: emptyCtx{},
		env:      env,
		eval:     noEvalPresent,
		receiver: receiver,
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

func WithReceiver(parent CallContext, receiver RubyObject) CallContext {
	if parent == nil {
		panic("parent CallContext cannot be nil")
	}
	return &wrappedCtx{parent, receiver}
}

type wrappedCtx struct {
	CallContext
	receiver RubyObject
}

func (c *wrappedCtx) Receiver() RubyObject { return c.receiver }

var (
	_ CallContext     = (*wrappedCtx)(nil)
	_ context.Context = (*wrappedCtx)(nil)
)

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

type callContext struct {
	emptyCtx
	env      Environment
	eval     func(node ast.Node, env Environment) (RubyObject, error)
	receiver RubyObject
}

func (c *callContext) Env() Environment { return c.env }
func (c *callContext) Eval(node ast.Node, env Environment) (RubyObject, error) {
	return c.eval(node, env)
}
func (c *callContext) Receiver() RubyObject { return c.receiver }

var (
	_ CallContext     = (*callContext)(nil)
	_ context.Context = (*callContext)(nil)
)
