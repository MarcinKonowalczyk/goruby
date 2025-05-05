package ast

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/utils"
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
	utils.AssertEqual(t, program.String(), "myVar = anotherVar")
}
