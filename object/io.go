package object

import (
	"bufio"
	"os"
)

var ioClass RubyClassObject = newClass(
	"IO",
	ioMethods,
	ioClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return NewIo(), nil
	},
)

func NewIo() *Io {
	return &Io{}
}

func init() {
	CLASSES.Set("Io", ioClass)
}

// Io represents a io in Ruby
type Io struct{}

var IO = &Io{}

// Inspect returns the Value
func (s *Io) Inspect() string { return "" }

// Type returns IO_OBJ
func (s *Io) Type() Type { return IO_OBJ }

// Class returns ioClass
func (s *Io) Class() RubyClass { return ioClass }

// // hashKey returns a hash key to be used by Hashes
// func (s *Io) hashKey() hashKey {
// 	h := fnv.New64a()
// 	h.Write([]byte(s.Value))
// 	return hashKey{Type: s.Type(), Value: h.Sum64()}
// }

var (
	_ RubyObject  = &Io{}
	_ inspectable = &Io{}
)

var ioClassMethods = map[string]RubyMethod{}

var ioMethods = map[string]RubyMethod{
	"gets": withArity(0, publicMethod(ioGets)),
}

func ioGets(context CallContext, args ...RubyObject) (RubyObject, error) {
	// read a string from stdin
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	// remove the newline character
	text = text[:len(text)-1]
	// create a new string object
	str := &String{Value: text}
	return str, nil
}
