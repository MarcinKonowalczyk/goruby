package parser_test

import (
	"fmt"

	"go/token"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/parser"
)

func ExampleParseFile() {
	fset := token.NewFileSet() // positions are relative to fset

	src := `LANG = "Ruby"

module Foo

	def bar()
		puts "Hello world"
	end

end`

	f, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the statements from the programs AST
	for _, s := range f.Statements {
		if exp, ok := s.(*ast.ExpressionStatement); ok {
			fmt.Printf("%T\n", exp.Expression)
		} else {
			fmt.Printf("%T\n", s)
		}
	}

	// output:
	//
	// *ast.Assignment
	// *ast.ModuleExpression
}

func ExampleParseExpr() {
	src := `def bar()
	puts "Hello world"
end`

	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%T\n", expr)

	// output:
	//
	// *ast.FunctionLiteral
}
