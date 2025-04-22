package object

func newMixin(class RubyClassObject, modules ...*Module) *mixin {
	return &mixin{class, modules}
}

type mixin struct {
	RubyClassObject
	modules []*Module
}

func (m *mixin) Methods() MethodSet {
	var methods = make(map[string]RubyMethod)
	for _, mod := range m.modules {
		moduleMethods := mod.Class().Methods()
		for k, v := range moduleMethods.GetAll() {
			methods[k] = v
		}
	}
	for k, v := range m.RubyClassObject.Methods().GetAll() {
		methods[k] = v
	}
	return NewMethodSet(methods)
}

func (m *mixin) GetMethod(name string) (RubyMethod, bool) {
	if method, ok := m.RubyClassObject.GetMethod(name); ok {
		return method, true
	}
	for _, mod := range m.modules {
		if method, ok := mod.Class().GetMethod(name); ok {
			return method, true
		}
	}
	return nil, false
}
