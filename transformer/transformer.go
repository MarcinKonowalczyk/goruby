package transformer

import (
	"fmt"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/MarcinKonowalczyk/goruby/transformer/units/block_lifts"
	"github.com/MarcinKonowalczyk/goruby/transformer/units/builtin_lifts"
)

type transformer struct {
	tracer trace.Tracer
}

func (t *transformer) Transform(node ast.Node) (ast.Node, error) {
	program, is_program := node.(*ast.Program)

	var transformed ast.Node
	var out_err error

	if is_program {
		transformed, out_err = t.transformProgram(program)

	} else {
		panic("Transforming non-program node is not yet implemented")
	}

	return transformed, out_err
}

func (t *transformer) transformProgram(program *ast.Program) (*ast.Program, error) {
	// var new_statement ast.Statement = context_call_expr_replacements[context_spec("find_all")].statement
	// program.Statements = append([]ast.Statement{new_statement}, program.Statements...)

	const ENABLE_BLOCK_LIFT_TRANSFORMER = true
	// const ENABLE_BLOCK_LIFT_TRANSFORMER = false
	if ENABLE_BLOCK_LIFT_TRANSFORMER {
		transformer := &block_lifts.Lift{}
		fmt.Printf("# TRANSFORMER applying %T pass 0\n", transformer)
		ast.Walk(program, transformer, nil)
		(*transformer).Pass = 1
		fmt.Printf("# TRANSFORMER applying %T pass 1\n", transformer)
		ast.Walk(program, transformer, nil)
		if len(transformer.LiftedFunctions) > 0 {
			var lifted []ast.Statement = make([]ast.Statement, len(transformer.LiftedFunctions))
			for _, function := range transformer.LiftedFunctions {
				lifted = append(lifted, &ast.ExpressionStatement{
					Expression: &ast.Assignment{
						Left:  &ast.Identifier{Value: function.Name},
						Right: function,
					},
				})
			}
			program.Statements = append(lifted, program.Statements...)
		}

		fmt.Printf("# TRANSFORMER done with %T\n", transformer)
	}

	const ENABLE_LIFT_TRANSFORMER = true
	// const ENABLE_LIFT_TRANSFORMER = false
	if ENABLE_LIFT_TRANSFORMER {
		var transformer = &builtin_lifts.Lift{}
		fmt.Printf("# TRANSFORMER applying %T\n", transformer)
		_ = ast.Walk(program, transformer, nil)
		if len(transformer.Statements) > 0 {
			program.Statements = append(transformer.Statements, program.Statements...)
		}
		fmt.Printf("# TRANSFORMER done with %T\n", transformer)
	}

	return program, nil
}
