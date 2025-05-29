package transformer

import (
	"context"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/walk"
	"github.com/MarcinKonowalczyk/goruby/trace"
	"github.com/MarcinKonowalczyk/goruby/transformer/logging"
	"github.com/MarcinKonowalczyk/goruby/transformer/units/function_lift"
)

type transformer struct {
	tracer trace.Tracer
}

func (t *transformer) Transform(node ast.Node, stages []Stage) (ast.Node, error) {
	return t.TransformCtx(context.Background(), node, stages)
}

func (t *transformer) TransformCtx(ctx context.Context, node ast.Node, stages []Stage) (ast.Node, error) {
	program, is_program := node.(*ast.Program)

	var transformed ast.Node
	var out_err error

	if is_program {
		transformed, out_err = t.transformProgram(ctx, program, stages)

	} else {
		panic("Transforming non-program node is not yet implemented")
	}

	return transformed, out_err
}

type Stage string

const (
	Sfunction_lift Stage = "function_lift"
)

var ALL_STAGES []Stage = []Stage{
	Sfunction_lift,
}

func (t *transformer) transformProgram(
	ctx context.Context,
	program *ast.Program,
	stages []Stage,
) (*ast.Program, error) {
	// var new_statement ast.Statement = context_call_expr_replacements[context_spec("find_all")].statement
	// program.Statements = append([]ast.Statement{new_statement}, program.Statements...)

	logging.Logf(ctx, "Transforming program with %d stages", len(stages))

	if len(stages) == 0 {
		return program, nil
	}

	var stages_map map[Stage]bool = make(map[Stage]bool)
	for _, stage := range stages {
		stages_map[stage] = true
	}

	// BLOCK_LIFT
	if _, ok := stages_map[Sfunction_lift]; ok {
		transformer := &function_lift.LiftBlocks{}
		logging.Logf(ctx, "=== applying %T pass 0 ===", transformer)
		walk.WalkCtx(ctx, program, transformer, nil)
		(*transformer).Pass = 1
		logging.Logf(ctx, "=== applying %T pass 1 ===", transformer)
		walk.WalkCtx(ctx, program, transformer, nil)
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

		logging.Logf(ctx, "=== done with %T ===", transformer)
	}

	if _, ok := stages_map[Sfunction_lift]; ok {
		var transformer = &function_lift.LiftBuiltins{}
		logging.Logf(ctx, "=== applying %T ===", transformer)
		_ = walk.WalkCtx(ctx, program, transformer, nil)
		if len(transformer.Statements) > 0 {
			program.Statements = append(transformer.Statements, program.Statements...)
		}
		logging.Logf(ctx, "=== done with %T ===", transformer)
	}

	return program, nil
}
