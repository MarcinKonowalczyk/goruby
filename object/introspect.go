package object

func getMethods(class RubyClass, addSuperMethods bool) *Array {
	var methodSymbols []RubyObject
	for class != nil {
		methods := class.Methods().GetAll()
		for meth, _ := range methods {
			methodSymbols = append(methodSymbols, &Symbol{meth})
		}
		if !addSuperMethods {
			break
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}
}
