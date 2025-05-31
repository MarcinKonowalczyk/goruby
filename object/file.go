package object

import (
	"os"
	"path/filepath"

	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

var fileClass ruby.ClassObject = newClass(
	"File",
	nil,
	fileClassMethods,
	notInstantiatable,
)

func init() {
	CLASSES.Set("File", fileClass)
}

var fileClassMethods = map[string]ruby.Method{
	"expand_path": newMethod(fileExpandPath),
	"dirname":     newMethod(fileDirname),
	"read":        newMethod(fileRead),
}

func fileExpandPath(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	switch len(args) {
	case 1:
		str, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(str, args[0])
		}
		path, err := filepath.Abs(str.Value)

		if err == nil {
			return NewString(path), nil
		}

		return nil, NewNotImplementedError("Cannot determine working directory")
	case 2:
		filename, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(filename, args[0])
		}
		dirname, ok := args[1].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(filename, args[0])
		}
		// TODO: make sure this is really the wanted behaviour
		abs, err := filepath.Abs(filepath.Join(dirname.Value, filename.Value))
		if err != nil {
			return nil, NewNotImplementedError("%s", err.Error())
		}

		return NewString(abs), nil
	default:
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
}

func fileDirname(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	filename, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(filename, args[0])
	}

	dirname := filepath.Dir(filename.Value)

	return NewString(dirname), nil
}

func fileRead(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, trace.HereCtx(ctx))()
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	filename, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(filename, args[0])
	}

	data, err := os.ReadFile(filename.Value)
	if err != nil {
		return nil, NewNotImplementedError("Cannot read file %s", filename.Value)
	}

	return NewString(string(data)), nil
}
