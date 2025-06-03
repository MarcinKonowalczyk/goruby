package evaluator

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/infix"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

// type callContext struct {
// 	object.CallContext
// 	evaluator Evaluator
// }

// func (c *callContext) Eval(node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error) {
// 	return c.evaluator.Eval(c, node, ev)
// }

// var (
// 	_ object.CallContext = &callContext{}
// )

type rubyObjects []ruby.Object

func (r rubyObjects) Inspect() string {
	toS := make([]string, len(r))
	for i, e := range r {
		toS[i] = e.Inspect()
	}
	return strings.Join(toS, ", ")
}
func (r rubyObjects) Class() ruby.Class { return nil }
func (r rubyObjects) HashKey() hash.Key { return expandToArrayIfNeeded(r).HashKey() }

func expandToArrayIfNeeded(obj ruby.Object) ruby.Object {
	arr, ok := obj.(rubyObjects)
	if !ok {
		return obj
	}
	return object.NewArray(arr...)
}

func Eval(ctx context.Context, node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	if node == nil {
		return nil, nil
	}
	if ev == nil {
		return nil, object.NewSyntaxError(fmt.Errorf("env is nil"))
	}

	e := NewEvaluator().SetOptions(func(opts *Options) {
	})

	res, err := e.Eval(ctx, node, ev)
	return res, err
}

type Evaluator interface {
	Eval(ctx context.Context, node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error)
	SetOptions(func(*Options)) Evaluator
}

type Options struct{}

func (e *evaluator) SetOptions(f func(*Options)) Evaluator {
	if f != nil {
		f(&e.options)
	}
	return e
}

func NewEvaluator() Evaluator {
	return &evaluator{}
}

var _ Evaluator = &evaluator{}

type evaluator struct {
	options Options // options for the evaluator
	ctx     call.Context[ruby.Object]
}

func (e *evaluator) Eval(ctx context.Context, node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	// e.ctx = &callContext{
	// 	CallContext: object.NewCallContext(env, object.FUNCS_STORE),
	// 	evaluator:   e,
	// }
	ectx := call.NewContext[ruby.Object](ctx, nil).
		WithReceiver(object.FUNCS_STORE).
		WithEnv(ev).
		WithEval(func(node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error) {
			return e.eval(node, ev)
		})
	e.ctx = ectx
	// ctx := object.NewCallContext(env, object.FUNCS_STORE)
	return e.eval(node, ev)
}

func (e *evaluator) eval(node ast.Node, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx, "evaluator.eval")()
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return e.evalProgram(node.Statements, ev)
	case *ast.ExpressionStatement:
		return e.evalExpressionStatement(node, ev)
	case *ast.ReturnStatement:
		return e.evalReturnStatement(node, ev)
	case *ast.BreakStatement:
		return e.evalBreakStatement(node, ev)
	case *ast.BlockStatement:
		return e.evalBlockStatement(node, ev)
	// Literals
	case *ast.IntegerLiteral:
		return e.evalIntegerLiteral(node, ev)
	case *ast.FloatLiteral:
		return e.evalFloatLiteral(node, ev)
	case *ast.Identifier:
		return e.evalIdentifier(node, ev)
	case *ast.StringLiteral:
		return e.evalStringLiteral(node, ev)
	case *ast.SymbolLiteral:
		return e.evalSymbolLiteral(node, ev)
	case *ast.FunctionLiteral:
		return e.evalFunctionLiteral(node, ev)
	case *ast.ArrayLiteral:
		return e.evalArrayLiteral(node, ev)
	case *ast.HashLiteral:
		return e.evalHashLiteral(node, ev)
	case ast.ExpressionList:
		return e.evalExpressionList(node, ev)
	// Expressions
	case *ast.Assignment:
		return e.evalAssignment(node, ev)
	case *ast.ContextCallExpression:
		return e.evalContextCallExpression(node, ev)
	case *ast.IndexExpression:
		return e.evalIndexExpression(node, ev)
	case *ast.PrefixExpression:
		right, err := e.eval(node.Right, ev)
		if err != nil {
			return nil, err
		}
		return e.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		return e.evalInfixExpression(node, ev)
	case *ast.ConditionalExpression:
		return e.evalConditionalExpression(node, ev)
	case *ast.Comment:
		trace.MessageCtx(e.ctx, "comment")
		// ignore comments
		return nil, nil
	case nil:
		trace.MessageCtx(e.ctx, "nil")
		return nil, nil
	case *ast.RangeLiteral:
		return e.evalRangeLiteral(node, ev)
	case *ast.Splat:
		return e.evalSplat(node, ev)
	case *ast.LoopExpression:
		return e.evalLoopExpression(node, ev)
	default:
		return nil, object.NewException("Unknown AST: %T", node)
	}

}

// NOTE: we support only double-quoted strings
var rep = []struct{ a, b string }{
	{a: "\\n", b: "\n"},
	{a: "\\t", b: "\t"},
	{a: "\\r", b: "\r"},
	{a: "\\b", b: "\b"},
	{a: "\\\\", b: "\\"},
}

func unescapeStringLiteral(node *ast.StringLiteral) string {
	value := node.Value
	for _, v := range rep {
		value = strings.ReplaceAll(value, v.a, v.b)
	}
	return value
}

var FORMAT_DIRECTIVE_RE = regexp.MustCompile(`\x60#\{(?P<content>[^}]*)\}\x60`)

func (e *evaluator) evalFormatDirectives(ev env.Environment[ruby.Object], value string) (string, error) {
	defer trace.TraceCtx(e.ctx)()
	if !strings.Contains(value, "`") {
		return value, nil
	}

	// example sting "hello `#{place}`"
	// search for `#{...}` pattern

	matches := FORMAT_DIRECTIVE_RE.FindAllStringSubmatchIndex(value, -1)
	if len(matches) == 0 {
		return value, nil
	}

	// for each match, evaluate the expression
	// note: in ruby this can be any expression but to save on parser evals
	// we only manually parse out the identifier name, and allow only that
	// hence `puts("hello `#{1+1}`")` *does* work in ruby but not here
	for _, match := range matches {
		start, end := match[0], match[1]
		start += 3 // skip the `#{`
		end -= 2   // skip the }`
		// loop up identifier
		val, err := e.eval(&ast.Identifier{Value: value[start:end]}, ev)
		if err != nil {
			return "", err
		}

		var val_str string
		switch val := val.(type) {
		case *object.String:
			val_str = val.Value
		case *object.Array:
			val_str = val.Inspect()
		default:
			fmt.Printf("val %T\n", val)
		}

		// replace the match with the value
		value = strings.Replace(value, value[match[0]:match[1]], val_str, 1)
	}

	return value, nil
}

func (e *evaluator) evalLoopExpression(node *ast.LoopExpression, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	for {
		value, err := e.evalBlockStatement(node.Block, ev)
		if err != nil {
			return nil, err
		}
		if value != nil {
			switch value := value.(type) {
			case *object.BreakValue:
				if object.IsTruthy(value.Value) {
					return value.Value, nil
				}
			case *object.ReturnValue:
				return value.Value, nil
			default:
				// <shrug> do nothing
			}
		}
	}
}

func (e *evaluator) evalProgram(statements []ast.Statement, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var result ruby.Object
	var err error
	for _, statement := range statements {
		if _, ok := statement.(*ast.Comment); ok {
			continue
		}
		result, err = e.eval(statement, ev)

		if err != nil {
			return nil, err
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}

	}
	return result, nil
}

func (e *evaluator) evalExpressions(expressions []ast.Expression, ev env.Environment[ruby.Object]) ([]ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var result []ruby.Object = make([]ruby.Object, 0, len(expressions))

	for _, expression := range expressions {
		evaluated, err := e.eval(expression, ev)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func (e *evaluator) evalArrayElements(elements []ast.Expression, ev env.Environment[ruby.Object]) ([]ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var result []ruby.Object

	for _, element := range elements {
		splat, ok := element.(*ast.Splat)
		if ok {
			// we're a splat! eval the splat and append the elements
			evaluated, err := e.eval(splat.Value, ev)
			if err != nil {
				return nil, err
			}
			if evaluated != nil {
				arrObj, ok := evaluated.(*object.Array)
				if !ok {
					return nil, object.NewException("splat value is not an array: %T", evaluated)
				}

				result = append(result, arrObj.Elements...)
			}
		} else {
			// not a splat. just eval the element
			evaluated, err := e.eval(element, ev)
			if err != nil {
				return nil, err
			}
			result = append(result, evaluated)
		}
	}
	return result, nil
}

func (e *evaluator) evalPrefixExpression(operator string, right ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch operator {
	case "!":
		return e.evalBangOperatorExpression(right), nil
	case "-":
		return e.evalMinusPrefixOperatorExpression(right)
	default:
		return nil, object.NewException("unknown operator: %s%s", operator, object.RubyObjectToTypeString(right))
	}
}

func (e *evaluator) evalBangOperatorExpression(right ruby.Object) ruby.Object {
	defer trace.TraceCtx(e.ctx)()
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NIL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func (e *evaluator) evalMinusPrefixOperatorExpression(right ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch right := right.(type) {
	case *object.Integer:
		return object.NewInteger(-right.Value), nil
	default:
		return nil, object.NewException("unknown operator: -%s", object.RubyObjectToTypeString(right))
	}
}

func (e *evaluator) evalConditionalExpression(ce *ast.ConditionalExpression, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	condition, err := e.eval(ce.Condition, ev)
	if err != nil {
		return nil, err
	}
	evaluateConsequence := object.IsTruthy(condition)
	if ce.Unless {
		evaluateConsequence = !evaluateConsequence
	}
	if evaluateConsequence {
		return e.eval(ce.Consequence, ev)
	} else if ce.Alternative != nil {
		return e.eval(ce.Alternative, ev)
	} else {
		return object.NIL, nil
	}
}

func (e *evaluator) evalIndexExpressionAssignment(left, index, right ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch target := left.(type) {
	case *object.Array:
		integer, ok := index.(*object.Integer)
		if !ok {
			return nil, object.NewImplicitConversionTypeError(integer, index)
		}
		idx := int(integer.Value)
		if idx >= len(target.Elements) {
			// enlarge slice
			for len(target.Elements) <= idx {
				target.Elements = append(target.Elements, object.NIL)
			}
		}
		target.Elements[idx] = right
		return right, nil
	case *object.Hash:
		target.Set(index, right)
		return right, nil
	default:
		return nil, object.NewException("assignment target not supported: %s", fmt.Sprintf("%T", left))
	}
}

func (e *evaluator) evalDefaultIndexExpression(left, index ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch target := left.(type) {
	case *object.Array:
		return e.evalArrayIndexExpression(target, index)
	case *object.Hash:
		return e.evalHashIndexExpression(target, index)
	case *object.String:
		return e.evalStringIndexExpression(target, index)
	default:
		return nil, object.NewException("index operator not supported: %s", fmt.Sprintf("%T", left))
	}
}

func (e *evaluator) evalSymbolIndexExpression(ev env.Environment[ruby.Object], target *object.Symbol, index ast.Expression) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch index.(type) {
	case *ast.Splat:
		trace.MessageCtx(e.ctx, "splat index expression")
		// evaluate the splat literal
		literal := index.(*ast.Splat)
		evaluated, err := e.eval(literal.Value, ev)
		if err != nil {
			return nil, err
		}
		if evaluated == nil {
			return nil, object.NewException("splat value is nil")
		}
		arrObj, ok := evaluated.(*object.Array)
		if !ok {
			return nil, object.NewException("splat value is not an array: %T", evaluated)
		}
		// call the proc with the splat arguments
		value, err := object.Send(e.ctx, target.Value, arrObj.Elements...)
		return value, err

	case ast.ExpressionList:
		trace.MessageCtx(e.ctx, "expression list index expression")
		// evaluate the expression list
		evaluated, err := e.evalExpressionList(index.(ast.ExpressionList), ev)
		if err != nil {
			return nil, err
		}
		if evaluated == nil {
			return nil, object.NewException("expression list is nil")
		}

		evaluated_obj := evaluated.(rubyObjects)

		// call the proc with the arguments
		args := make([]ruby.Object, len(evaluated_obj))
		for i, e := range evaluated_obj {
			if e == nil {
				args[i] = object.NIL
			} else {
				args[i] = e
			}
		}
		value, err := object.Send(e.ctx, target.Value, args...)
		if err != nil {
			trace.MessageCtx(e.ctx, "error sending to proc: %s", err)
		}
		return value, err

	default:
		trace.MessageCtx(e.ctx, "default index expression")
		evaluated, err := e.eval(index, ev)
		if err != nil {
			return nil, err
		}

		if evaluated == nil {
			return nil, object.NewException("symbol index is nil")
		}

		// call the proc with the arguments
		args := make([]ruby.Object, 1)
		if evaluated == nil {
			args[0] = object.NIL
		} else {
			args[0] = evaluated
		}

		value, err := object.Send(e.ctx, target.Value, args...)
		return value, err
	}
}

func (e *evaluator) evalArrayIndexExpression(arrayObject *object.Array, index ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	len := int64(len(arrayObject.Elements))
	switch index := index.(type) {
	case *object.Integer:
		idx, out_of_bounds := objectIntegerToIndex(index, len)
		if out_of_bounds {
			return object.NIL, nil
		}
		return arrayObject.Elements[idx], nil
	case *object.Array:
		left, length, out_of_bounds, err := objectArrayToIndex(index, len)
		if err != nil {
			return nil, err
		}
		if out_of_bounds {
			return object.NIL, nil
		}
		return object.NewArray(arrayObject.Elements[left:(left + length)]...), nil
	case *object.Range:
		left, right, out_of_bounds, err := objectRangeToIndex(index, len)
		if err != nil {
			return nil, err
		}
		if out_of_bounds {
			return object.NIL, nil
		}
		return object.NewArray(arrayObject.Elements[left:right]...), nil
	case rubyObjects:
		// we got a bunch of objects as the index
		index_array := object.NewArray(index...)
		return e.evalArrayIndexExpression(arrayObject, index_array)
	default:
		err := &object.TypeError{
			Message: fmt.Sprintf("array index must be Integer, Array or Range, got %T", index),
		}

		return nil, err
	}
}

func (e *evaluator) evalHashIndexExpression(hash *object.Hash, index ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	result, ok := hash.Get(index)
	if !ok {
		return object.NIL, nil
	}
	return result, nil
}

func int64ToIndex(idx int64, len int64) (int64, bool) {
	max_negative := -len
	max_positive := max_negative*-1 - 1
	if max_positive < 0 {
		return 0, true
	}
	if idx > 0 && idx > max_positive {
		return 0, true
	}
	if idx < 0 && idx < max_negative {
		return 0, true
	}
	// wrap negative index
	for idx < 0 {
		idx += len
	}
	if idx < 0 {
		return 0, true
	}
	return idx, false
}

func objectIntegerToIndex(index *object.Integer, len int64) (int64, bool) {
	return int64ToIndex(index.Value, len)
}

func objectArrayToIndex(index *object.Array, _ int64) (int64, int64, bool, error) {
	if len(index.Elements) == 0 {
		return 0, 0, true, nil
	}
	first_element := index.Elements[0]
	last_element := index.Elements[len(index.Elements)-1]

	left_index, ok := first_element.(*object.Integer)
	if !ok {
		return 0, 0, true, object.NewImplicitConversionTypeError(first_element, first_element)
	}
	left_idx := left_index.Value

	length_index, ok := last_element.(*object.Integer)
	if !ok {
		return 0, 0, true, object.NewImplicitConversionTypeError(last_element, last_element)
	}
	length_idx := length_index.Value

	return left_idx, length_idx, false, nil
}

func objectRangeToIndex(index *object.Range, length int64) (int64, int64, bool, error) {
	left_idx, out_of_bounds := int64ToIndex(index.Left, length)
	if out_of_bounds {
		return 0, 0, true, nil
	}
	right_idx, out_of_bounds := int64ToIndex(index.Right, length)
	if out_of_bounds {
		return 0, 0, true, nil
	}

	if index.Inclusive {
		right_idx++
	}

	return left_idx, right_idx, false, nil
}

func (e *evaluator) evalStringIndexExpression(stringObject *object.String, index ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	switch index := index.(type) {
	case *object.Integer:
		idx, out_of_bounds := objectIntegerToIndex(index, int64(len(stringObject.Value)))
		if out_of_bounds {
			return object.NIL, nil
		}
		return object.NewString(string(stringObject.Value[idx])), nil
	case *object.Array:
		left, length, out_if_bounds, err := objectArrayToIndex(index, int64(len(stringObject.Value)))
		if err != nil {
			return nil, err
		}
		if out_if_bounds {
			return object.NIL, nil
		}
		return object.NewString(string(stringObject.Value[left:(left + length)])), nil
	case *object.Range:
		left, right, out_if_bounds, err := objectRangeToIndex(index, int64(len(stringObject.Value)))
		if err != nil {
			return nil, err
		}
		if out_if_bounds {
			return object.NIL, nil
		}
		return object.NewString(string(stringObject.Value[left:right])), nil
	case rubyObjects:
		// we got a bunch of objects as the index
		index_array := object.NewArray(index...)
		return e.evalStringIndexExpression(stringObject, index_array)

	default:
		// fmt.Printf("index: %s(%T)\n", index.Inspect(), index)
		err := object.NewImplicitConversionTypeErrorMany(index, object.NewInteger(0), object.NewFloat(0.0))
		return nil, err
	}

}

func (e *evaluator) evalBlockStatement(block *ast.BlockStatement, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var result ruby.Object
	var err error
	for _, statement := range block.Statements {
		result, err = e.eval(statement, ev)
		if err != nil {
			return nil, err
		}
		switch result.(type) {
		case nil:
			// do nothing
		case *object.ReturnValue:
			return result, nil
		case *object.BreakValue:
			if object.IsTruthy(result.(*object.BreakValue).Value) {
				return result, nil
			}
		}
	}
	if result == nil {
		return object.NIL, nil
	}
	return result, nil
}

func (e *evaluator) evalIdentifier(node *ast.Identifier, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	trace.MessageCtx(e.ctx, node.Value)

	val, ok := ev.Get(node.Value)
	if ok {
		return val, nil
	}

	if node.IsConstant() {
		return nil, object.NewUninitializedConstantNameError(node.Value)
	}

	if node.IsGlobal() {
		return object.NIL, nil
	}

	// maybe a function
	// fmt.Println("ident", node)
	val, err := object.Send(e.ctx, node.Value)
	if err != nil {
		return nil, object.NewNoMethodError(object.FUNCS_STORE, node.Value)
	}
	return val, nil
}

func (e *evaluator) evalExpressionStatement(node *ast.ExpressionStatement, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	return e.eval(node.Expression, ev)
}

func (e *evaluator) evalReturnStatement(node *ast.ReturnStatement, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	val, err := e.eval(node.ReturnValue, ev)
	if err != nil {
		return nil, err
	}
	return &object.ReturnValue{Value: val}, nil
}

func (e *evaluator) evalBreakStatement(node *ast.BreakStatement, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var val ruby.Object
	val, err := e.eval(node.Condition, ev)
	if err != nil {
		return nil, err
	}
	if node.Unless {
		if object.IsTruthy(val) {
			val = object.FALSE
		} else {
			val = object.TRUE
		}
	}
	return &object.BreakValue{Value: val}, nil
}

func (e *evaluator) evalStringLiteral(node *ast.StringLiteral, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	value := unescapeStringLiteral(node)
	value, err := e.evalFormatDirectives(ev, value)
	if err != nil {
		return nil, err
	}
	return object.NewString(value), nil
}

func (e *evaluator) evalFunctionLiteral(node *ast.FunctionLiteral, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	// context, _ := env.Get("bottom")
	// construct a function object and stick it onto self
	params := make([]*object.FunctionParameter, len(node.Parameters))
	for i, param := range node.Parameters {
		def, err := e.eval(param.Default, ev)
		if err != nil {
			return nil, err
		}
		params[i] = &object.FunctionParameter{Name: param.Name, Default: def, Splat: param.Splat}
	}
	function := &object.Function{
		Name:       node.Name,
		Parameters: params,
		Env:        ev,
		Body:       node.Body,
	}
	_, extended := ruby.AddMethod(object.FUNCS_STORE, node.Name, function)
	if extended {
		panic("we should not be extending FUNCS. they already should be extended")
	}
	// if extended {
	// 	// we've just extended the context. set it in the env. this should not normally fire
	// 	env.Set("bottom", newContext)
	// }
	return object.NewSymbol(node.Name), nil
}

func (e *evaluator) evalArrayLiteral(node *ast.ArrayLiteral, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	elements, err := e.evalArrayElements(node.Elements, ev)
	if err != nil {
		return nil, err
	}
	// TODO: If any of the elements is a splat, we need to flatten them
	return object.NewArray(elements...), nil
}

func (e *evaluator) evalHashLiteral(node *ast.HashLiteral, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var hash object.Hash
	for k, v := range node.Map {
		key, err := e.eval(k, ev)
		if err != nil {
			return nil, err
		}
		value, err := e.eval(v, ev)
		if err != nil {
			return nil, err
		}
		hash.Set(key, value)
	}
	return &hash, nil
}

func (e *evaluator) evalExpressionList(node ast.ExpressionList, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	var objects []ruby.Object
	for _, n := range node {
		obj, err := e.eval(n, ev)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return rubyObjects(objects), nil

}

func (e *evaluator) evalAssignment(node *ast.Assignment, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	right, err := e.eval(node.Right, ev)
	if err != nil {
		return nil, err
	}

	switch left := node.Left.(type) {
	case *ast.IndexExpression:
		indexLeft, err := e.eval(left.Left, ev)
		if err != nil {
			return nil, err
		}
		index, err := e.eval(left.Index, ev)
		if err != nil {
			return nil, err
		}
		return e.evalIndexExpressionAssignment(indexLeft, index, expandToArrayIfNeeded(right))
	case *ast.Identifier:
		right = expandToArrayIfNeeded(right)
		if left.IsGlobal() {
			ev.SetGlobal(left.Value, right)
		} else {
			ev.Set(left.Value, right)
		}
		return right, nil
	case ast.ExpressionList:
		var values rubyObjects
		switch right := right.(type) {
		case rubyObjects:
			values = right
		case *object.Array:
			values = right.Elements
		default:
			values = []ruby.Object{right}
		}
		if len(left) > len(values) {
			// enlarge slice
			for len(values) <= len(left) {
				values = append(values, object.NIL)
			}
		}
		for i, exp := range left {
			switch exp := exp.(type) {
			case *ast.Identifier:
				ev.Set(exp.Value, values[i])
			case *ast.Splat:
				panic("splat in assignment not implemented yet")
			case *ast.IndexExpression:
				indexLeft, err := e.eval(exp.Left, ev)
				if err != nil {
					return nil, err
				}
				index, err := e.eval(exp.Index, ev)
				if err != nil {
					return nil, err
				}
				_, err = e.evalIndexExpressionAssignment(indexLeft, index, values[i])
				if err != nil {
					return nil, err
				}
				continue
			default:
				panic("unexpected expression in assignment: " + fmt.Sprintf("%T", exp))
			}
		}
		return expandToArrayIfNeeded(right), nil
	default:
		return nil, object.NewSyntaxError(fmt.Errorf("assignment not supported to %T", node.Left))
	}
}

func (e *evaluator) evalContextCallExpression(node *ast.ContextCallExpression, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	trace.MessageCtx(e.ctx, node.Function)

	context, err := e.eval(node.Context, ev)
	if err != nil {
		return nil, err
	}
	if context == nil {
		context = object.FUNCS_STORE
	}

	args, err := e.evalExpressions(node.Arguments, ev)
	if err != nil {
		return nil, err
	}
	if node.Block != nil {
		block, err := e.eval(node.Block, ev)
		if err != nil {
			return nil, err
		}
		args = append(args, block)
	}
	return object.Send(e.ctx.WithReceiver(context), node.Function, args...)
}

func (e *evaluator) evalIndexExpression(node *ast.IndexExpression, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	left, err := e.eval(node.Left, ev)
	if err != nil {
		return nil, err
	}
	switch left := left.(type) {
	case *object.Symbol:
		// functions evaluate to symbols with the name of the function
		// anonymous functions evaluate to a functions with a random name
		// indexing them should call them
		// NOTE: we pass unevaluated index to proc
		return e.evalSymbolIndexExpression(ev, left, node.Index)
	default:
		index, err := e.eval(node.Index, ev)
		if err != nil {
			return nil, err
		}
		return e.evalDefaultIndexExpression(left, index)
	}
}

func (e *evaluator) evalInfixExpression(node *ast.InfixExpression, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	left, err := e.eval(node.Left, ev)
	if err != nil {
		return nil, err
	}

	if node.Operator == infix.LOGICALOR {
		if object.IsTruthy(left) {
			// left is already truthy. don't evaluate right side
			return left, nil
		}
	} else if node.Operator == infix.LOGICALAND {
		if !object.IsTruthy(left) {
			// left is already falsy. don't evaluate right side
			return left, nil
		}
	}

	right, err := e.eval(node.Right, ev)
	if err != nil {
		return nil, err
	}

	if node.Operator == infix.LOGICALOR {
		// left is not truthy, since we're here
		// result is right
		return right, nil
	} else if node.Operator == infix.LOGICALAND {
		// left is not falsy, since we're here
		// result is right
		return right, nil
	}
	return object.Send(e.ctx.WithReceiver(left), node.Operator.String(), right)
}

func (e *evaluator) evalRangeLiteral(node *ast.RangeLiteral, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()

	left, err := e.eval(node.Left, ev)
	if err != nil {
		return nil, err
	}
	right, err := e.eval(node.Right, ev)
	if err != nil {
		return nil, err
	}
	if left == nil || right == nil {
		return nil, object.NewSyntaxError(fmt.Errorf("range start or end is nil"))
	}
	leftInt, ok := left.(*object.Integer)
	if !ok {
		return nil, object.NewSyntaxError(fmt.Errorf("range start is not an integer: %T", left))
	}
	rightInt, ok := right.(*object.Integer)
	if !ok {
		return nil, object.NewSyntaxError(fmt.Errorf("range end is not an integer: %T", right))
	}
	return &object.Range{
		Left:      leftInt.Value,
		Right:     rightInt.Value,
		Inclusive: node.Inclusive,
	}, nil
}

func (e *evaluator) evalSplat(node *ast.Splat, ev env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	val, err := e.eval(node.Value, ev)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, object.NewSyntaxError(fmt.Errorf("splat value is nil"))
	}

	switch val := val.(type) {
	case *object.Array:
		return object.NewArray(val.Elements...), nil
	default:
		return object.NewArray(val), nil
	}
}

func (e *evaluator) evalIntegerLiteral(node *ast.IntegerLiteral, _ env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	return object.NewInteger(node.Value), nil
}

func (e *evaluator) evalFloatLiteral(node *ast.FloatLiteral, _ env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	return object.NewFloat(node.Value), nil
}

func (e *evaluator) evalSymbolLiteral(node *ast.SymbolLiteral, _ env.Environment[ruby.Object]) (ruby.Object, error) {
	defer trace.TraceCtx(e.ctx)()
	return object.NewSymbol(node.Value), nil
}
