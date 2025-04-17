package object

var (
	nilClass RubyClassObject = newClass(
		"NilClass", objectClass, nilMethods, nilClassMethods, notInstantiatable,
	)
	// NIL represents the singleton object nil
	NIL RubyObject = &nilObject{}
)

func init() {
	classes.Set("NilClass", nilClass)
}

type nilObject struct{}

func (n *nilObject) Inspect() string  { return "nil" }
func (n *nilObject) Type() Type       { return NIL_OBJ }
func (n *nilObject) Class() RubyClass { return nilClass }

var nilClassMethods = map[string]RubyMethod{}

var nilMethods = map[string]RubyMethod{
	"nil?": withArity(0, publicMethod(nilIsNil)),
	// "==":   withArity(1, publicMethod(nilEqual)),
}

func nilIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return TRUE, nil
}

// func nilEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
// 	if len(args) == 0 {
// 		return nil, NewArgumentError("wrong number of arguments (given 0, expected 1)")
// 	}
// 	if args[0] == NIL {
// 		return TRUE, nil
// 	}
// 	return FALSE, nil
// }
