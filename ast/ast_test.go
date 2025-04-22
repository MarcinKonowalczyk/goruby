package ast

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &Assignment{
					Left: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Right: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}
	if program.String() != "myVar = anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
