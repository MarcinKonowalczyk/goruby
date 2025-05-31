package object

import (
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

func WithArity(arity int, fn ruby.Method) ruby.Method {
	return ruby.NewMethod(
		func(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
			defer trace.TraceCtx(ctx, "withArity")()
			if len(args) != arity {
				return nil, NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(ctx, args...)
		},
	)
}
