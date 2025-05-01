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

var objectClass = newClass(
	"Object",
	nil,
	objectMethodSet,
	objectClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &Object{}, nil
	})

func init() {
	classes.Set("Object", objectClass)
}

// Object represents an Object in Ruby
type Object struct {
	_ int // for uniqueness
}

// Inspect return ""
func (o *Object) Inspect() string { return "" }

// Type returns OBJECT_OBJ
func (o *Object) Type() Type { return OBJECT_OBJ }

// Class returns objectClass
func (o *Object) Class() RubyClass { return objectClass }

var objectClassMethods = map[string]RubyMethod{}

var objectMethodSet = map[string]RubyMethod{
	"to_s":    withArity(0, publicMethod(objectToS)),
	"is_a?":   withArity(1, publicMethod(objectIsA)),
	"nil?":    withArity(0, publicMethod(objectIsNil)),
	"methods": publicMethod(objectMethods),
	"class":   withArity(0, publicMethod(objectClassMethod)),
	"puts":    publicMethod(objectPuts),
	"print":   publicMethod(objectPrint),
	"require": withArity(1, publicMethod(objectRequire)),
	"tap":     publicMethod(objectTap),
	"raise":   publicMethod(objectRaise),
	"==":      withArity(1, publicMethod(objectEqual)),
	"!=":      withArity(1, publicMethod(objectNotEqual)),
}

func objectToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if self, ok := receiver.(*Self); ok {
		receiver = self.RubyObject
	}
	val := fmt.Sprintf("#<%s:%p>", receiver.Class().Name(), receiver)
	return &String{Value: val}, nil
}

func objectIsA(context CallContext, args ...RubyObject) (RubyObject, error) {
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

func print(lines []string, delimiter string) {
	var out strings.Builder
	for i, line := range lines {
		out.WriteString(line)
		if i != len(lines)-1 {
			out.WriteString(delimiter)
		}
	}
	out.WriteString(delimiter)
	fmt.Print(out.String())
}

func objectPuts(context CallContext, args ...RubyObject) (RubyObject, error) {
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
				if arg == NIL {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	print(lines, "\n")
	return NIL, nil
}

func objectPrint(context CallContext, args ...RubyObject) (RubyObject, error) {
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
				if arg == NIL {
					//
				} else {
					lines = append(lines, arg.Inspect())
				}
			default:
				lines = append(lines, arg.Inspect())
			}
		}
	}
	print(lines, "")
	return NIL, nil
}

func objectMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperMethods := true
	if len(args) == 1 {
		if val, ok := SymbolToBool(args[0]); ok {
			showSuperMethods = val
		}
	}

	receiver := context.Receiver()
	class := context.Receiver().Class()

	extended, ok := receiver.(*extendedObject)

	if !showSuperMethods && !ok {
		return &Array{}, nil
	}

	if !showSuperMethods && ok {
		class = extended.class
	}

	return getMethods(class, showSuperMethods), nil
}

func objectIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if receiver == NIL {
		return TRUE, nil
	}
	return FALSE, nil
}

func objectClassMethod(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if _, ok := receiver.(RubyClassObject); ok {
		return classClass, nil
	}
	return receiver.Class().(RubyClassObject), nil
}

func objectRequire(context CallContext, args ...RubyObject) (RubyObject, error) {
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

func objectTap(context CallContext, args ...RubyObject) (RubyObject, error) {
	block := args[0]
	proc, ok := block.(*Symbol)
	if !ok {
		return nil, NewArgumentError("map requires a block")
	}
	self, _ := context.Env().Get("self")
	self_class := self.Class()
	fn, ok := self_class.GetMethod(proc.Value)
	if !ok {
		return nil, NewNoMethodError(self, proc.Value)
	}
	_, err := fn.Call(context, context.Receiver())
	if err != nil {
		return nil, err
	}
	return context.Receiver(), nil
}

func objectRaise(context CallContext, args ...RubyObject) (RubyObject, error) {
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

func swapOrFalse(left, right RubyObject, swapped bool) bool {
	if swapped {
		// we've already swapped. just return false
		return false
	} else {
		// 1-depth recursive call with swapped arguments
		return rubyObjectsEqual(right, left, true)
	}
}

func rubyObjectsEqual(left, right RubyObject, swapped bool) bool {
	// leftClass := left.Class()
	// rightClass := right.Class()
	// if leftClass != rightClass {
	// 	return swapOrFalse(left, right, swapped)
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
	case *Integer:
		right_t, ok := safeObjectToInteger(right)
		if !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *Float:
		right_t, ok := safeObjectToFloat(right)
		if !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t
		}
	case *String:
		if right_t, ok := right.(*String); !ok {
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	case *Array:
		if right_t, ok := right.(*Array); !ok {
			return swapOrFalse(left, right, swapped)
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
			return swapOrFalse(left, right, swapped)
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
			return swapOrFalse(left, right, swapped)
		} else {
			return left.Value == right_t.Value
		}
	default:
		return false
	}
}

func RubyObjectsEqual(left, right RubyObject) bool {
	return rubyObjectsEqual(left, right, false)
}

func objectEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	if RubyObjectsEqual(context.Receiver(), args[0]) {
		return TRUE, nil
	}
	return FALSE, nil
}

func objectNotEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	if RubyObjectsEqual(context.Receiver(), args[0]) {
		return FALSE, nil
	}
	return TRUE, nil
}
