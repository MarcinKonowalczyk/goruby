package object

import (
	"fmt"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/parser"
	"github.com/pkg/errors"
)

var kernelModule = newModule("Kernel", kernelMethodSet, nil)

func init() {
	classes.Set("Kernel", kernelModule)
}

var kernelMethodSet = map[string]RubyMethod{
	"to_s":              withArity(0, publicMethod(kernelToS)),
	"is_a?":             withArity(1, publicMethod(kernelIsA)),
	"nil?":              withArity(0, publicMethod(kernelIsNil)),
	"methods":           publicMethod(kernelMethods),
	"public_methods":    publicMethod(kernelPublicMethods),
	"protected_methods": publicMethod(kernelProtectedMethods),
	"private_methods":   publicMethod(kernelPrivateMethods),
	"class":             withArity(0, publicMethod(kernelClass)),
	"puts":              privateMethod(kernelPuts),
	"print":             privateMethod(kernelPrint),
	"require":           withArity(1, privateMethod(kernelRequire)),
	"extend":            publicMethod(kernelExtend),
	"block_given?":      withArity(0, privateMethod(kernelBlockGiven)),
	"tap":               publicMethod(kernelTap),
	"raise":             privateMethod(kernelRaise),
	"==":                withArity(1, publicMethod(kernelEqual)),
	"!=":                withArity(1, publicMethod(kernelNotEqual)),
}

func kernelToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if self, ok := receiver.(*Self); ok {
		receiver = self.RubyObject
	}
	val := fmt.Sprintf("#<%s:%p>", receiver.Class().Name(), receiver)
	return &String{Value: val}, nil
}

func kernelIsA(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	receiver_class := context.Receiver().Class()
	switch arg := args[0].(type) {
	case RubyClassObject:
		if arg.Name() == receiver_class.Name() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return nil, NewTypeError("argument must be a Class")
	}
}

func argsToLines(args []RubyObject) []string {
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
			case *nilObject:
				//
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	return lines
}

func print(lines []string, end bool) {
	var out strings.Builder
	for i, line := range lines {
		out.WriteString(line)
		if i != len(lines)-1 {
			out.WriteString("\n")
		}
	}
	if end {
		out.WriteString("\n")
	}
	fmt.Print(out.String())
}

func kernelPuts(context CallContext, args ...RubyObject) (RubyObject, error) {
	lines := argsToLines(args)
	print(lines, true)
	return NIL, nil
}

func kernelPrint(context CallContext, args ...RubyObject) (RubyObject, error) {
	lines := argsToLines(args)
	print(lines, false)
	return NIL, nil
}

func kernelMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showInstanceMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showInstanceMethods = boolean.Value
	}

	receiver := context.Receiver()
	class := context.Receiver().Class()

	extended, ok := receiver.(*extendedObject)

	if !showInstanceMethods && !ok {
		return &Array{}, nil
	}

	if !showInstanceMethods && ok {
		class = extended.class
	}

	publicMethods := getMethods(class, PUBLIC_METHOD, showInstanceMethods)
	protectedMethods := getMethods(class, PROTECTED_METHOD, showInstanceMethods)
	return &Array{Elements: append(publicMethods.Elements, protectedMethods.Elements...)}, nil
}

func kernelPublicMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PUBLIC_METHOD, showSuperClassMethods), nil
}

func kernelProtectedMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PROTECTED_METHOD, showSuperClassMethods), nil
}

func kernelPrivateMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PRIVATE_METHOD, showSuperClassMethods), nil
}

func kernelIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return FALSE, nil
}

func kernelClass(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if _, ok := receiver.(RubyClassObject); ok {
		return classClass, nil
	}
	return receiver.Class().(RubyClassObject), nil
}

func kernelRequire(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	name, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(name, args[0])
	}
	filename := name.Value
	if !strings.HasSuffix(filename, "rb") {
		filename += ".rb"
	}
	absolutePath, _ := filepath.Abs(filename)
	loadedFeatures, ok := context.Env().Get("$LOADED_FEATURES")
	if !ok {
		loadedFeatures = NewArray()
		context.Env().SetGlobal("$LOADED_FEATURES", loadedFeatures)
	}
	arr, ok := loadedFeatures.(*Array)
	if !ok {
		arr = NewArray()
	}
	loaded := false
	for _, feat := range arr.Elements {
		if feat.Inspect() == absolutePath {
			loaded = true
			break
		}
	}
	if loaded {
		return FALSE, nil
	}

	file, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		found := false
		loadPath, _ := context.Env().Get("$:")
		for _, p := range loadPath.(*Array).Elements {
			newPath := path.Join(p.(*String).Value, filename)
			file, err = os.ReadFile(newPath)
			if !os.IsNotExist(err) {
				absolutePath = newPath
				found = true
				break
			}
		}
		if !found {
			return nil, NewNoSuchFileLoadError(name.Value)
		}
	}

	prog, err := parser.ParseFile(token.NewFileSet(), absolutePath, file, 0)
	if err != nil {
		return nil, NewSyntaxError(err)
	}
	_, err = context.Eval(prog, WithScopedLocalVariables(context.Env()))
	if err != nil {
		return nil, errors.WithMessage(err, "require")
	}
	arr.Elements = append(arr.Elements, &String{Value: absolutePath})
	return TRUE, nil
}

func kernelExtend(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) == 0 {
		return nil, NewWrongNumberOfArgumentsError(1, 0)
	}
	modules := make([]*Module, len(args))
	for i, arg := range args {
		module, ok := arg.(*Module)
		if !ok {
			return nil, NewWrongArgumentTypeError(module, arg)
		}
		modules[i] = module
	}
	extended := &extendedObject{
		RubyObject: context.Receiver(),
		class: newEigenclass(
			newMixin(context.Receiver().Class().(RubyClassObject), modules...),
			map[string]RubyMethod{},
		),
	}
	info, _ := EnvStat(context.Env(), context.Receiver())
	info.Env().Set(info.Name(), extended)
	return extended, nil
}

func kernelBlockGiven(context CallContext, args ...RubyObject) (RubyObject, error) {
	self, _ := context.Receiver().(*Self)
	if self.Block == nil {
		return FALSE, nil
	}
	return TRUE, nil
}

func kernelTap(context CallContext, args ...RubyObject) (RubyObject, error) {
	block, remainingArgs, ok := extractBlockFromArgs(args)
	if !ok {
		return nil, NewNoBlockGivenLocalJumpError()
	}
	if len(remainingArgs) != 0 {
		return nil, NewWrongNumberOfArgumentsError(0, 1)
	}
	_, err := block.Call(context, context.Receiver())
	if err != nil {
		return nil, err
	}
	return context.Receiver(), nil
}

func kernelRaise(context CallContext, args ...RubyObject) (RubyObject, error) {
	switch len(args) {
	case 1:
		switch arg := args[0].(type) {
		case *String:
			return nil, NewRuntimeError("%s", arg.Value)
		default:
			exc, err := Send(NewCallContext(context.Env(), arg), "exception")
			if err != nil {
				return nil, NewTypeError("exception class/object expected")
			}
			if excAsErr, ok := exc.(error); ok {
				return nil, excAsErr
			}
			return nil, nil
		}
	default:
		return nil, NewRuntimeError("")
	}
}

func swapOr(result bool, left, right RubyObject, swapped bool) bool {
	if swapped {
		// we've already swapped. just return the result
		return result
	} else {
		// 1-depth recursive call with swapped arguments
		return rubyObjectsEqual(right, left, true)
	}
}

func rubyObjectsEqual(left, right RubyObject, swapped bool) bool {
	// leftClass := left.Class()
	// rightClass := right.Class()
	// if leftClass != rightClass {
	// 	return swapOr(false, left, right, swapped)
	// }
	// if left == nil {
	// 	return right == nil || right.Class().Name() == "NilClass"
	// }
	// if left.Class().Name() == "NilClass" {
	// 	return right == nil || right.Class().Name() == "NilClass"
	// }
	// fmt.Printf("left: %T right: %T\n", left, right)
	// fmt.Println("left:", left.Class().Name(), "right:", right.Class().Name())
	switch left := left.(type) {
	case *Boolean:
		if right_t, ok := right.(*Boolean); !ok {
			// swap. maybe the other thing knows how to compare
			// itself to a boolean
			return swapOr(false, left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	case *Integer:
		right_t, ok := safeObjectToInteger(right)
		if !ok {
			return swapOr(false, left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *Float:
		right_t, ok := safeObjectToFloat(right)
		if !ok {
			return swapOr(false, left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *String:
		if right_t, ok := right.(*String); !ok {
			return swapOr(false, left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	case *Array:
		if right_t, ok := right.(*Array); !ok {
			return swapOr(false, left, right, swapped)
		} else {
			if len(left.Elements) != len(right_t.Elements) {
				return false
			}
			for i, elem := range left.Elements {
				if !rubyObjectsEqual(elem, right_t.Elements[i], swapped) {
					return false
				}
			}
			return true
		}
	case *Hash:
		if right_t, ok := right.(*Hash); !ok {
			return swapOr(false, left, right, swapped)
		} else {
			if len(left.Map) != len(right_t.Map) {
				return false
			}
			for key, leftValue := range left.ObjectMap() {
				rightValue, ok := right_t.Get(key)
				if !ok {
					return false
				}
				if !rubyObjectsEqual(leftValue, rightValue, swapped) {
					return false
				}
			}

			return true
		}
	case *Symbol:
		if right_t, ok := right.(*Symbol); !ok {
			return swapOr(false, left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	case *nilObject:
		if right_t, ok := right.(*nilObject); !ok {
			return swapOr(false, left, right, swapped)
		} else {
			return left == right_t
		}
	default:
		return false
	}
}

func RubyObjectsEqual(left, right RubyObject) bool {
	return rubyObjectsEqual(left, right, false)
}

func kernelEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	// if self, ok := receiver.(*Self); ok {
	// 	receiver = self.RubyObject
	// }
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	arg := args[0]
	res := RubyObjectsEqual(receiver, arg)
	if res {
		return TRUE, nil
	}
	return FALSE, nil
}

func kernelNotEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	// if self, ok := receiver.(*Self); ok {
	// 	receiver = self.RubyObject
	// }
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	arg := args[0]
	res := RubyObjectsEqual(receiver, arg)
	if res {
		return FALSE, nil
	}
	return TRUE, nil
}
