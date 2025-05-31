package object

import (
	"bufio"
	"os"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var IoClass ruby.ClassObject = newClass("IO", nil, ioClassMethods)

func init() {
	CLASSES.Set("Io", IoClass)
}

var ioClassMethods = map[string]ruby.Method{
	"gets": WithArity(0, ruby.NewMethod(ioClassGets)),
}

func ioClassGets(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
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
