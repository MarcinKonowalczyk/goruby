package object

func getMethods(class RubyClass, addSuperMethods bool) *Array {
	var methodSymbols []RubyObject
	methods := class.Methods().GetAll()
	for meth := range methods {
		methodSymbols = append(methodSymbols, &Symbol{meth})
	}
	if class == bottomClass {
		return &Array{Elements: methodSymbols}
	}
	if addSuperMethods {
		methods := bottomClass.Methods().GetAll()
		for meth := range methods {
			methodSymbols = append(methodSymbols, &Symbol{meth})
		}
	}
	return &Array{Elements: methodSymbols}
}
