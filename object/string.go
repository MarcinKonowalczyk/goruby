package object

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
)

var stringClass RubyClassObject = newClass(
	"String",
	objectClass,
	stringMethods,
	stringClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &String{}, nil
	},
)

func init() {
	classes.Set("String", stringClass)
}

// String represents a string in Ruby
type String struct {
	Value string
}

// Inspect returns the Value
func (s *String) Inspect() string { return s.Value }

// Type returns STRING_OBJ
func (s *String) Type() Type { return STRING_OBJ }

// Class returns stringClass
func (s *String) Class() RubyClass { return stringClass }

// hashKey returns a hash key to be used by Hashes
func (s *String) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hashKey{Type: s.Type(), Value: h.Sum64()}
}

var (
	_ RubyObject  = &String{}
	_ inspectable = &String{}
)

func stringify(obj RubyObject) (string, error) {
	stringObj, err := Send(NewCallContext(nil, obj), "to_s")
	if err != nil {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s into String",
				obj.Class().Name(),
			),
		)
	}
	str, ok := stringObj.(*String)
	if !ok {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s to String (%s#to_s gives %s)",
				obj.Class().Name(),
				obj.Class().Name(),
				stringObj.Class().Name(),
			),
		)
	}
	return str.Value, nil
}

var stringClassMethods = map[string]RubyMethod{}

var stringMethods = map[string]RubyMethod{
	"initialize": privateMethod(stringInitialize),
	"to_s":       withArity(0, publicMethod(stringToS)),
	"+":          withArity(1, publicMethod(stringAdd)),
	"gsub":       withArity(2, publicMethod(stringGsub)),
	"length":     withArity(0, publicMethod(stringLength)),
	"size":       withArity(0, publicMethod(stringLength)),
	// "==":         withArity(1, publicMethod(stringEqual)),
	// "!=":         withArity(1, publicMethod(stringNotEqual)),
	"lines": withArity(0, publicMethod(stringLines)),
	"to_f":  withArity(0, publicMethod(stringToF)),
}

func stringInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	self, _ := context.Receiver().(*Self)
	switch len(args) {
	case 0:
		self.RubyObject = &String{}
		return self, nil
	case 1:
		str, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(str, args[0])
		}
		self.RubyObject = &String{Value: str.Value}
		return self, nil
	default:
		return nil, NewWrongNumberOfArgumentsError(len(args), 1)
	}
}

func stringToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	str := context.Receiver().(*String)
	return &String{str.Value}, nil
}

func stringAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	add, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(add, args[0])
	}
	return &String{s.Value + add.Value}, nil
}

func stringGsub(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	pattern, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(pattern, args[0])
	}
	replacement, ok := args[1].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(replacement, args[1])
	}

	// Perform the gsub operation
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, NewTypeError(fmt.Sprintf("Invalid regex pattern: %s", err))
	}

	result := re.ReplaceAllString(s.Value, replacement.Value)

	// Return the modified string
	return &String{Value: result}, nil
}

func stringLength(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	return &Integer{Value: int64(len(s.Value))}, nil
}

// func stringEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	s := context.Receiver().(*String)
// 	other, ok := args[0].(*String)
// 	if !ok {
// 		return nil, NewImplicitConversionTypeError(other, args[0])
// 	}
// 	return &Boolean{Value: s.Value == other.Value}, nil
// }

// func stringNotEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	s := context.Receiver().(*String)
// 	other, ok := args[0].(*String)
// 	if !ok {
// 		return nil, NewImplicitConversionTypeError(other, args[0])
// 	}
// 	return &Boolean{Value: s.Value != other.Value}, nil
// }

func stringLines(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	lines := strings.Split(s.Value, "\n")
	arr := NewArray()
	for _, line := range lines {
		arr.Elements = append(arr.Elements, &String{Value: line + "\n"})
		// arr.Elements = append(arr.Elements, &String{Value: line})
	}
	return arr, nil
}

var FLOAT_RE = regexp.MustCompile(`[-+]?\d*\.?\d+`)

func stringToF(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	if s.Value == "" {
		return &Float{Value: 0.0}, nil
	}
	match := FLOAT_RE.FindString(s.Value)
	if match == "" {
		return &Float{Value: 0.0}, nil
	}
	// Convert the string to a float
	val, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return nil, NewTypeError("Invalid float value: " + s.Value)
	}
	return &Float{Value: val}, nil
}
