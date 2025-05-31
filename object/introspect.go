package object

import "github.com/MarcinKonowalczyk/goruby/object/ruby"

func getMethods(cls ruby.Class, addSuperMethods bool) *Array {
	var methodSymbols []ruby.Object
	names := cls.Methods().Names()
	for _, name := range names {
		methodSymbols = append(methodSymbols, NewSymbol(name))
	}
	if c, ok := cls.(*class); ok && c == bottomClass {
		return &Array{Elements: methodSymbols}
	}
	if addSuperMethods {
		names := bottomClass.Methods().Names()
		for _, name := range names {
			methodSymbols = append(methodSymbols, NewSymbol(name))
		}
	}
	return &Array{Elements: methodSymbols}
}
