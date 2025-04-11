package object

import (
	"hash/fnv"
)

var regexClass RubyClassObject = newClass(
	"Regex",
	objectClass,
	regexMethods,
	regexClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &Regex{}, nil
	},
)

func init() {
	classes.Set("Regex", regexClass)
}

// Regex represents a regex in Ruby
type Regex struct {
	Value     string
	Modifiers string
}

// Inspect returns the Value
func (s *Regex) Inspect() string { return "/" + s.Value + "/" + s.Modifiers }

// Type returns REGEX_OBJ
func (s *Regex) Type() Type { return REGEX_OBJ }

// Class returns regexClass
func (s *Regex) Class() RubyClass { return regexClass }

// hashKey returns a hash key to be used by Hashes
func (s *Regex) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hashKey{Type: s.Type(), Value: h.Sum64()}
}

var (
	_ RubyObject  = &Regex{}
	_ inspectable = &Regex{}
)

var regexClassMethods = map[string]RubyMethod{}

var regexMethods = map[string]RubyMethod{
	"initialize": privateMethod(regexInitialize),
	"to_s":       withArity(0, publicMethod(regexToS)),
}

func regexInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	self, _ := context.Receiver().(*Self)
	switch len(args) {
	case 0:
		self.RubyObject = &Regex{}
		return self, nil
	case 1:
		reg, ok := args[0].(*Regex)
		if !ok {
			return nil, NewImplicitConversionTypeError(reg, args[0])
		}
		self.RubyObject = &Regex{Value: reg.Value}
		return self, nil
	default:
		return nil, NewWrongNumberOfArgumentsError(len(args), 1)
	}
}

func regexToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	reg := context.Receiver().(*Regex)
	return &String{reg.Inspect()}, nil
}
