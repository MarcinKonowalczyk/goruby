package object

import (
	"bufio"
	"os"

	"github.com/MarcinKonowalczyk/goruby/trace"
)

var IoClass RubyClassObject = newClass(
	"IO",
	nil,
	ioClassMethods,
	notInstantiatable,
)

func init() {
	CLASSES.Set("Io", IoClass)
}

var ioClassMethods = map[string]RubyMethod{
	"gets": withArity(0, newMethod(ioClassGets)),
}

func ioClassGets(ctx CallContext, args ...RubyObject) (RubyObject, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	// read a string from stdin
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// remove the newline character
	text = text[:len(text)-1]
	// create a new string object
	str := NewString(text)
	return str, nil
}
