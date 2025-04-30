package evaluator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/ast/infix"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/pkg/errors"
)

type callContext struct {
	object.CallContext
}

func (c *callContext) Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	return Eval(node, env)
}

type rubyObjects []object.RubyObject

func (r rubyObjects) Inspect() string {
	toS := make([]string, len(r))
	for i, e := range r {
		toS[i] = e.Inspect()
	}
	return strings.Join(toS, ", ")
}
func (r rubyObjects) Type() object.Type       { return "" }
func (r rubyObjects) Class() object.RubyClass { return nil }

func expandToArrayIfNeeded(obj object.RubyObject) object.RubyObject {
	arr, ok := obj.(rubyObjects)
	if !ok {
		return obj
	}
	return object.NewArray(arr...)
}

func my_debug_panic_(msg string) {
	// fmt.Println("my_debug_panic", msg)
	// panic(msg)
}

// Eval evaluates the given node and traverses recursive over its children
func Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	// fmt.Println("Eval", node, fmt.Sprintf("%T", node))
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval of return statement")
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.BreakStatement:
		var val object.RubyObject
		val, err := Eval(node.Condition, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval of break statement")
		}
		if node.Unless {
			if isTruthy(val) {
				val = object.FALSE
			} else {
				val = object.TRUE
			}
		}
		return &object.BreakValue{Value: val}, nil
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Literals
	case (*ast.IntegerLiteral):
		return object.NewInteger(node.Value), nil
	case (*ast.FloatLiteral):
		return object.NewFloat(node.Value), nil
	// case (*ast.Boolean):
	// 	return nativeBoolToBooleanObject(node.Value), nil
	case (*ast.Keyword__FILE__):
		my_debug_panic_("case (*ast.Keyword__FILE__):")
		return &object.String{Value: node.Filename}, nil
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Global:
		val, ok := env.Get(node.Value)
		if !ok {
			return object.NIL, nil
		}
		return val, nil
	case *ast.StringLiteral:
		value := unescapeStringLiteral(node)
		value, err := evaluateFormatDirectives(env, value)
		if err != nil {
			return nil, errors.WithMessage(err, "eval string literal")
		}
		return &object.String{Value: value}, nil
	case *ast.SymbolLiteral:
		switch node.Value {
		case "true":
			return object.TRUE, nil
		case "false":
			return object.FALSE, nil
		case "nil":
			return object.NIL, nil
		default:
			return &object.Symbol{Value: node.Value}, nil
		}
	case *ast.FunctionLiteral:
		context, _ := env.Get("self")
		params := make([]*object.FunctionParameter, len(node.Parameters))
		for i, param := range node.Parameters {
			def, err := Eval(param.Default, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function literal param")
			}
			params[i] = &object.FunctionParameter{Name: param.Name.Value, Default: def, Splat: param.Splat}
		}
		function := &object.Function{
			Parameters: params,
			Env:        env,
			Body:       node.Body,
		}
		object.AddMethod(context, node.Name.Value, function)
		return &object.Symbol{Value: node.Name.Value}, nil

	case *ast.ArrayLiteral:
		elements, err := evalArrayElements(node.Elements, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval array literal")
		}
		// TODO: If any of the elements is a splat, we need to flatten them
		return &object.Array{Elements: elements}, nil
	case *ast.HashLiteral:
		var hash object.Hash
		for k, v := range node.Map {
			key, err := Eval(k, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash key")
			}
			value, err := Eval(v, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash value")
			}
			hash.Set(key, value)
		}
		return &hash, nil
	case ast.ExpressionList:
		my_debug_panic_("case ast.ExpressionList:")
		var objects []object.RubyObject
		for _, e := range node {
			obj, err := Eval(e, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval expression list")
			}
			objects = append(objects, obj)
		}
		return rubyObjects(objects), nil

	// Expressions
	case *ast.Assignment:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval right hand Assignment side")
		}

		switch left := node.Left.(type) {
		case *ast.IndexExpression:
			indexLeft, err := Eval(left.Left, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
			}
			index, err := Eval(left.Index, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
			}
			return evalIndexExpressionAssignment(indexLeft, index, expandToArrayIfNeeded(right))
		case *ast.Identifier:
			right = expandToArrayIfNeeded(right)
			env.Set(left.Value, right)
			return right, nil
		case *ast.Global:
			right = expandToArrayIfNeeded(right)
			env.SetGlobal(left.Value, right)
			return right, nil
		case ast.ExpressionList:
			var values rubyObjects
			switch right := right.(type) {
			case rubyObjects:
				my_debug_panic_("case rubyObjects:")
				values = right
			case *object.Array:
				values = right.Elements
			default:
				my_debug_panic_("case default: (2)")
				values = []object.RubyObject{right}
			}
			if len(left) > len(values) {
				// enlarge slice
				for len(values) <= len(left) {
					values = append(values, object.NIL)
				}
			}
			for i, exp := range left {
				if indexExp, ok := exp.(*ast.IndexExpression); ok {
					indexLeft, err := Eval(indexExp.Left, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
					}
					index, err := Eval(indexExp.Index, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
					}
					_, err = evalIndexExpressionAssignment(indexLeft, index, values[i])
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
					}
					continue
				}
				env.Set(exp.String(), values[i])
			}
			return expandToArrayIfNeeded(right), nil
		default:
			my_debug_panic_("case default: (3)")
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("assignment not supported to %T", node.Left)),
			)
		}
	case *ast.ContextCallExpression:
		context, err := Eval(node.Context, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call receiver")
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call arguments")
		}
		if node.Block != nil {
			block, err := Eval(node.Block, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval method call block")
			}
			args = append(args, block)
		}
		callContext := &callContext{object.NewCallContext(env, context)}
		return object.Send(callContext, node.Function.Value, args...)
	case *ast.IndexExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval IndexExpression left side")
		}
		switch left := left.(type) {
		case *object.Symbol:
			// functions evaluate to symbols with the name of the function
			// anonymous functions evaluate to a functions with a random name
			// indexing them should call them
			// NOTE: we pass unevaluated index to proc
			return evalSymbolIndexExpression(env, left, node.Index)
		default:
			index, err := Eval(node.Index, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval IndexExpression index")
			}
			return evalIndexExpression(left, index)
		}
	case *ast.PrefixExpression:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval prefix right side")
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator left side")
		}

		if node.Operator == infix.LOGICALOR {
			if isTruthy(left) {
				// left is altready truthy. don't evaluate right side
				return left, nil
			}
		} else if node.Operator == infix.LOGICALAND {
			if !isTruthy(left) {
				// left is altready falsy. don't evaluate right side
				return left, nil
			}
		}

		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator right side")
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
		context := &callContext{object.NewCallContext(env, left)}
		return object.Send(context, node.Operator.String(), right)

	case *ast.ConditionalExpression:
		return evalConditionalExpression(node, env)

	case *ast.Comment:
		// ignore comments
		return nil, nil

	case nil:
		return nil, nil

	case *ast.RangeLiteral:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval range start")
		}
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval range end")
		}
		if left == nil || right == nil {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start or end is nil")),
			)
		}
		if left.Type() != right.Type() {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start and end are not the same type: %s %s", left.Type(), right.Type())),
			)
		}
		if left.Type() != object.INTEGER_OBJ {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start and end are not integers: %s %s", left.Type(), right.Type())),
			)
		}
		leftInt, ok := left.(*object.Integer)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start is not an integer: %s", left.Type())),
			)
		}
		rightInt, ok := right.(*object.Integer)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range end is not an integer: %s", right.Type())),
			)
		}

		// left_int := leftInt.Value
		// right_int := rightInt.Value

		// flip := false
		// if left_int > right_int {
		// 	flip = true
		// 	left_int, right_int = right_int, left_int
		// }

		// if node.Inclusive {
		// 	right_int++
		// }

		// elements := make(rubyObjects, right_int-left_int)
		// for i := left_int; i < right_int; i++ {
		// 	elements[i-left_int] = &object.Integer{Value: i}
		// }

		// if flip {
		// 	// reverse the elements
		// 	slices.Reverse(elements)
		// }

		return &object.Range{
			Left:      leftInt,
			Right:     rightInt,
			Inclusive: node.Inclusive,
		}, nil

	case *ast.Splat:
		my_debug_panic_("case *ast.Splat:")

		val, err := Eval(node.Value, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval splat value")
		}
		if val == nil {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("splat value is nil")),
			)
		}

		switch val := val.(type) {
		case *object.Array:
			my_debug_panic_("case *object.Array:")
			return &object.Array{
				Elements: val.Elements,
			}, nil
		default:
			my_debug_panic_("case default: (5)")
			return &object.Array{
				Elements: []object.RubyObject{val},
			}, nil
		}

	case *ast.LoopExpression:

		return evalLoopExpression(node, env)

	default:
		my_debug_panic_("case default: (6)")
		err := object.NewException("Unknown AST: %T", node)
		return nil, errors.WithStack(err)
	}

}

func unescapeStringLiteral(node *ast.StringLiteral) string {
	rep := map[string]string{
		"\\n":  "\n",
		"\\t":  "\t",
		"\\r":  "\r",
		"\\b":  "\b",
		"\\\\": "\\",
	}
	// NOTE: we support only double-quoted strings
	rep["\""] = "\""
	value := node.Value
	for k, v := range rep {
		value = strings.ReplaceAll(value, k, v)
	}
	return value
}

var FORMAT_DIRECTIVE_RE = regexp.MustCompile(`\x60#\{(?P<content>[^}]*)\}\x60`)

func evaluateFormatDirectives(env object.Environment, value string) (string, error) {
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
		val, err := Eval(&ast.Identifier{
			Value: value[start:end],
		}, env)
		if err != nil {
			return "", errors.WithMessage(err, "eval format directive")
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

func evalLoopExpression(node *ast.LoopExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(node.Condition, env)
	if err != nil {
		return nil, errors.WithMessage(err, "eval loop condition")
	}
	for isTruthy(condition) {
		value, err := evalBlockStatement(node.Block, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval loop body")
		}
		if value != nil {
			switch value := value.(type) {
			case *object.BreakValue:
				if isTruthy(value.Value) {
					return value.Value, nil
				}
			case *object.ReturnValue:
				return value.Value, nil
			default:
				// <shrug> do nothing
			}
		}
		condition, err = Eval(node.Condition, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval loop condition")
		}
	}
	return object.NIL, nil
}

func evalProgram(stmts []ast.Statement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range stmts {
		if _, ok := statement.(*ast.Comment); ok {
			continue
		}
		result, err = Eval(statement, env)

		if err != nil {
			return nil, errors.WithMessage(err, "eval program statement")
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}

	}
	return result, nil
}

func evalExpressions(exps []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalArrayElements(elements []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range elements {
		splat, ok := e.(*ast.Splat)
		if ok {
			// we're a splat! eval the splat and append the elements
			evaluated, err := Eval(splat.Value, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval splat value")
			}
			if evaluated != nil {
				if evaluated.Type() != object.ARRAY_OBJ {
					return nil, errors.WithStack(
						object.NewException("splat value is not an array: %s", evaluated.Type()),
					)
				}
				arrObj, ok := evaluated.(*object.Array)
				if !ok {
					return nil, errors.WithStack(
						object.NewException("splat value is not an array: %s", evaluated.Type()),
					)
				}

				result = append(result, arrObj.Elements...)
			}
		} else {
			// not a splat. just eval the element
			evaluated, err := Eval(e, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval array elements")
			}
			result = append(result, evaluated)
		}
	}
	return result, nil
}

func evalPrefixExpression(operator string, right object.RubyObject) (object.RubyObject, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: %s%s", operator, right.Type()))
	}
}

func evalBangOperatorExpression(right object.RubyObject) object.RubyObject {
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

func evalMinusPrefixOperatorExpression(right object.RubyObject) (object.RubyObject, error) {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}, nil
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: -%s", right.Type()))
	}
}

func evalConditionalExpression(ce *ast.ConditionalExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(ce.Condition, env)
	if err != nil {
		return nil, err
	}
	evaluateConsequence := isTruthy(condition)
	if ce.Unless {
		evaluateConsequence = !evaluateConsequence
	}
	if evaluateConsequence {
		return Eval(ce.Consequence, env)
	} else if ce.Alternative != nil {
		return Eval(ce.Alternative, env)
	} else {
		return object.NIL, nil
	}
}

func evalIndexExpressionAssignment(left, index, right object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		integer, ok := index.(*object.Integer)
		if !ok {
			return nil, errors.Wrap(
				object.NewImplicitConversionTypeError(integer, index),
				"eval array index",
			)
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
		return nil, errors.Wrap(
			object.NewException("assignment target not supported: %s", left.Type()),
			"eval IndexExpression Assignment",
		)
	}
}

func evalIndexExpression(left, index object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		return evalArrayIndexExpression(target, index)
	case *object.Hash:
		return evalHashIndexExpression(target, index)
	case *object.String:
		return evalStringIndexExpression(target, index)
	default:
		var left_type string = string(left.Type())
		if left_type == "" {
			left_type = fmt.Sprintf("%T", left)
		}
		return nil, errors.WithStack(object.NewException("index operator not supported: %s", left_type))
	}
}

func evalSymbolIndexExpression(env object.Environment, target *object.Symbol, index ast.Expression) (object.RubyObject, error) {
	switch index.(type) {
	case *ast.Splat:
		// evaluate the splat literal
		literal := index.(*ast.Splat)
		evaluated, err := Eval(literal.Value, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval splat value")
		}
		if evaluated == nil {
			return nil, errors.WithStack(
				object.NewException("splat value is nil"),
			)
		}
		if evaluated.Type() != object.ARRAY_OBJ {
			return nil, errors.WithStack(
				object.NewException("splat value is not an array: %s", evaluated.Type()),
			)
		}
		arrObj, ok := evaluated.(*object.Array)
		if !ok {
			return nil, errors.WithStack(
				object.NewException("splat value is not an array: %s", evaluated.Type()),
			)
		}

		// call the proc with the splat arguments
		args := make([]object.RubyObject, len(arrObj.Elements))
		for i, e := range arrObj.Elements {
			if e == nil {
				args[i] = object.NIL
			} else {
				args[i] = e
			}
		}
		printable_args := make([]string, len(args))
		for i, e := range args {
			if e == nil {
				printable_args[i] = "nil"
			} else {
				printable_args[i] = e.Inspect()
			}
		}
		self, _ := env.Get("self")
		callContext := &callContext{object.NewCallContext(env, self)}
		// value, err := target.Call(callContext, args...)

		// get the method from the env
		// method, ok := env.Get(target.Value)
		// if !ok {
		// 	return nil, errors.WithStack(
		// 		object.NewException("method not found: %s", target.Value),
		// 	)
		// }
		// call the method
		// fmt.Println(target.Value, "(", strings.Join(printable_args, ", "), ")")
		value, err := object.Send(callContext, target.Value, args...)
		return value, err
	default:
		// not implemented yet
		return nil, errors.WithStack(
			object.NewException("proc index operator not supported: %s", index),
		)
	}
}

func evalArrayIndexExpression(arrayObject *object.Array, index object.RubyObject) (object.RubyObject, error) {
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
			return nil, errors.Wrap(err, "array index array")
		}
		if out_of_bounds {
			return object.NIL, nil
		}
		return &object.Array{Elements: arrayObject.Elements[left:(left + length)]}, nil
	case *object.Range:
		left, right, out_of_bounds, err := objectRangeToIndex(index, len)
		if err != nil {
			return nil, errors.Wrap(err, "array index range")
		}
		if out_of_bounds {
			return object.NIL, nil
		}
		return &object.Array{Elements: arrayObject.Elements[left:right]}, nil
	case rubyObjects:
		// we got a bunch of objects as the index
		index_array := object.NewArray(index...)
		return evalArrayIndexExpression(arrayObject, index_array)
	default:
		index_type := string(index.Type())
		if index_type == "" {
			index_type = fmt.Sprintf("%T", index)
		}
		err := &object.TypeError{
			Message: fmt.Sprintf("array index must be Integer, Array or Range, got %s", index_type),
		}

		return nil, err
	}
}

func evalHashIndexExpression(hash *object.Hash, index object.RubyObject) (object.RubyObject, error) {
	result, ok := hash.Get(index)
	if !ok {
		return object.NIL, nil
	}
	return result, nil
}

func objectIntegerToIndex(index *object.Integer, len int64) (int64, bool) {
	idx := index.Value
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

func objectArrayToIndex(index *object.Array, length int64) (int64, int64, bool, error) {
	if len(index.Elements) == 0 {
		return 0, 0, true, nil
	}
	first_element := index.Elements[0]
	last_element := index.Elements[len(index.Elements)-1]

	if first_element.Type() != object.INTEGER_OBJ {
		return 0, 0, true, errors.Wrap(
			object.NewImplicitConversionTypeError(first_element, first_element),
			"eval array index",
		)
	}

	if last_element.Type() != object.INTEGER_OBJ {
		return 0, 0, true, errors.Wrap(
			object.NewImplicitConversionTypeError(last_element, last_element),
			"eval array index",
		)
	}

	left_index, ok := first_element.(*object.Integer)
	if !ok {
		return 0, 0, true, errors.Wrap(
			object.NewImplicitConversionTypeError(first_element, first_element),
			"eval array index",
		)
	}
	left_idx := left_index.Value

	length_index, ok := last_element.(*object.Integer)
	if !ok {
		return 0, 0, true, errors.Wrap(
			object.NewImplicitConversionTypeError(last_element, last_element),
			"eval array index",
		)
	}
	length_idx := length_index.Value

	return left_idx, length_idx, false, nil
}

func objectRangeToIndex(index *object.Range, length int64) (int64, int64, bool, error) {
	left := index.Left
	right := index.Right

	left_idx, out_of_bounds := objectIntegerToIndex(left, length)
	if out_of_bounds {
		return 0, 0, true, nil
	}
	right_idx, out_of_bounds := objectIntegerToIndex(right, length)
	if out_of_bounds {
		return 0, 0, true, nil
	}

	if index.Inclusive {
		right_idx++
	}

	return left_idx, right_idx, false, nil
}

func evalStringIndexExpression(stringObject *object.String, index object.RubyObject) (object.RubyObject, error) {
	switch index := index.(type) {
	case *object.Integer:
		idx, out_of_bounds := objectIntegerToIndex(index, int64(len(stringObject.Value)))
		if out_of_bounds {
			return object.NIL, nil
		}
		return &object.String{Value: string(stringObject.Value[idx])}, nil
	case *object.Array:
		left, length, out_if_bounds, err := objectArrayToIndex(index, int64(len(stringObject.Value)))
		if err != nil {
			return nil, errors.Wrap(err, "string index array")
		}
		if out_if_bounds {
			return object.NIL, nil
		}
		return &object.String{Value: string(stringObject.Value[left:(left + length)])}, nil
	case *object.Range:
		left, right, out_if_bounds, err := objectRangeToIndex(index, int64(len(stringObject.Value)))
		if err != nil {
			return nil, errors.Wrap(err, "string index range")
		}
		if out_if_bounds {
			return object.NIL, nil
		}
		return &object.String{Value: string(stringObject.Value[left:right])}, nil
	case rubyObjects:
		// we got a bunch of objects as the index
		index_array := object.NewArray(index...)
		return evalStringIndexExpression(stringObject, index_array)

	default:
		// fmt.Printf("index: %s(%T)\n", index.Inspect(), index)
		err := errors.Wrap(
			object.NewImplicitConversionTypeErrorMany(index, object.NewInteger(0), object.NewFloat(0.0)),
			"eval string index",
		)
		return nil, err
	}

}

func evalBlockStatement(block *ast.BlockStatement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range block.Statements {
		result, err = Eval(statement, env)
		if err != nil {
			return nil, err
		}
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ {
				return result, nil
			} else if rt == object.BREAK_VALUE_OBJ {
				if isTruthy(result.(*object.BreakValue).Value) {
					return result, nil
				} else {
					//
				}
			}
		}
	}
	if result == nil {
		return object.NIL, nil
	}
	return result, nil
}

func evalIdentifier(node *ast.Identifier, env object.Environment) (object.RubyObject, error) {
	val, ok := env.Get(node.Value)
	if ok {
		return val, nil
	}

	if node.Constant {
		return nil, errors.Wrap(
			object.NewUninitializedConstantNameError(node.Value),
			"eval identifier",
		)
	}

	self, _ := env.Get("self")
	context := &callContext{object.NewCallContext(env, self)}
	val, err := object.Send(context, node.Value)
	if err != nil {
		return nil, errors.Wrap(
			object.NewUndefinedLocalVariableOrMethodNameError(self, node.Value),
			"eval ident as method call",
		)
	}
	return val, nil
}

//	func unwrapReturnValue(obj object.RubyObject) object.RubyObject {
//		if returnValue, ok := obj.(*object.ReturnValue); ok {
//			return returnValue.Value
//		}
//		return obj
//	}
func isTruthy(obj object.RubyObject) bool {
	switch obj {
	case object.NIL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Symbol:
			val, ok := object.SymbolToBool(obj)
			if ok {
				// Special boolean symbols are their respective values
				return val
			} else {
				// Other symbols are truthy
				return true
			}
		case *object.Integer:
			return obj.Value != 0
		case *object.Float:
			return obj.Value != 0.0
		case *object.String:
			return obj.Value != ""
		case *object.Array:
			return len(obj.Elements) > 0
		case *object.Hash:
			return len(obj.Map) > 0
		default:
			return true
		}
	}
}

// IsError returns true if the given RubyObject is an object.Error or an
// object.Exception (or any subclass of object.Exception)
func IsError(obj object.RubyObject) bool {
	if obj != nil {
		return obj.Type() == object.EXCEPTION_OBJ
	}
	return false
}
