package object

import (
	"bufio"
	"io"
	"reflect"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

type ioClass struct {
	class
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewIoClass(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) ruby.ClassObject {
	ioClass := &ioClass{*ioClassBase, stdin, stdout, stderr}
	CLASSES.Set("Io", ioClass)
	return ioClass
}

var ioClassBase = newClass("Io", nil, ioClassMethods)

var ioClassMethods = map[string]ruby.Method{
	"gets":  WithArity(0, ruby.NewMethod(ioClassGets)),
	"puts":  WithArity(1, ruby.NewMethod(ioPuts)),
	"print": WithArity(1, ruby.NewMethod(ioPrint)),
}

func ioClassGets(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	self, ok := ctx.Receiver().(*ioClass)
	if !ok {
		return nil, NewTypeError("expected io class receiver")
	}
	if self.stdin == nil {
		return nil, NewRuntimeError("no stdin available")
	}
	// read a string from stdin
	reader := bufio.NewReader(self.stdin)
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

func isNil(thing any) bool {
	if thing == nil {
		return true
	}
	typeof := reflect.TypeOf(thing)
	if typeof.Kind() == reflect.Ptr {
		if thing == reflect.Zero(typeof).Interface() {
			return true
		}
	}
	return false
}

func print(stdout io.Writer, lines []string, delimiter string) error {
	if isNil(stdout) {
		return nil
	}
	var out strings.Builder
	for i, line := range lines {
		out.WriteString(line)
		if i != len(lines)-1 {
			out.WriteString(delimiter)
		}
	}
	out.WriteString(delimiter)
	_, err := stdout.Write([]byte(out.String()))
	return err
}

func ioPuts(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	self, ok := ctx.Receiver().(*ioClass)
	if !ok {
		return nil, NewTypeError("expected io class receiver")
	}
	var lines []string
	for _, arg := range args {
		if arr, ok := arg.(*Array); ok {
			// arg is an array. splat it out
			// todo: make it a deep splat? check with original ruby implementation
			for _, elem := range arr.Elements {
				lines = append(lines, elem.Inspect())
			}
		} else {
			switch arg := arg.(type) {
			case *Symbol:
				if arg == NIL.(*Symbol) {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	err := print(self.stdout, lines, "\n")
	return NIL, err
}

func ioPrint(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx)()
	self, ok := ctx.Receiver().(*ioClass)
	if !ok {
		return nil, NewTypeError("expected io class receiver")
	}
	var lines []string
	for _, arg := range args {
		if arr, ok := arg.(*Array); ok {
			// arg is an array. splat it out
			// todo: make it a deep splat? check with original ruby implementation
			// for _, elem := range arr.Elements {
			// 	lines = append(lines, elem.Inspect())
			// }
			lines = append(lines, arr.Inspect())
		} else {
			switch arg := arg.(type) {
			case *Symbol:
				if arg == NIL.(*Symbol) {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	err := print(self.stdout, lines, "")
	return NIL, err
}
