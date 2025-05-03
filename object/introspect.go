package object

func getMethods(class RubyClass, addSuperMethods bool) *Array {
	var methodSymbols []RubyObject
	names := class.Methods().Names()
	for _, name := range names {
		methodSymbols = append(methodSymbols, &Symbol{name})
	}
	if class == bottomClass {
		return &Array{Elements: methodSymbols}
	}
	if addSuperMethods {
		names := bottomClass.Methods().Names()
		for _, name := range names {
			methodSymbols = append(methodSymbols, &Symbol{name})
		}
	}
	return &Array{Elements: methodSymbols}
}
