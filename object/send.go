package object

import (
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

// Send sends message method with args to context and returns its result
func Send(ctx call.Context[ruby.Object], method string, args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	trace.MessageCtx(ctx, method)

	receiver := ctx.Receiver()
	receiver_class := receiver.Class()

	// search for the method in the ancestry tree
	for receiver_class != nil {
		fn, ok := receiver_class.GetMethod(method)
		if !ok {
			if rc, ok := receiver_class.(*class); ok && rc == bottomClass {
				// no method and we are at the top of the ancestry tree
				break
			}
			receiver_class = bottomClass
			continue
		}

		return fn.Call(ctx, args...)
	}

	// fmt.Printf("receiver: %v(%T)\n", receiver, receiver)
	// fmt.Printf("method: %v(%T)\n", method, method)
	return nil, NewNoMethodError(receiver, method)
}
