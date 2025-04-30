package ast

import (
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &Assignment{
					Left: &Identifier{
						Value: "myVar",
					},
					Right: &Identifier{
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
